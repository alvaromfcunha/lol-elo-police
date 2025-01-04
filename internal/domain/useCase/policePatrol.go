package usecase

import (
	"errors"
	"os"
	"slices"

	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity/enum"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/repository"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/service"
)

type PolicePatrol struct {
	PlayerRepository           repository.IPlayerRepository
	RankedInfoRepository       repository.IRankedInfoRepository
	MatchRepository            repository.IMatchRepository
	MatchParticipantRepository repository.IMatchParticipantRepository
	LolService                 service.ILolService
	WhatsappService            service.IWhatsappService
	TemplateService            service.ITemplateService
}

func (u PolicePatrol) Execute() error {
	players, err := u.PlayerRepository.GetAll()
	if err != nil {
		return errors.New("cannot get all players on database")
	}

	type checkPlayerHasNewMatch struct {
		Player     entity.Player
		IsNewMatch bool
		MatchId    string
		Err        error
	}

	checkPlayersChannel := make(chan checkPlayerHasNewMatch, len(players))
	for _, player := range players {
		go func(player entity.Player, channel chan checkPlayerHasNewMatch) {
			isNewMatch, matchId, err := u.checkPlayerHasNewMatch(player)
			channel <- checkPlayerHasNewMatch{player, isNewMatch, matchId, err}
		}(player, checkPlayersChannel)
	}

	results := make([]checkPlayerHasNewMatch, len(players))
	for i := range players {
		results[i] = <-checkPlayersChannel
		err = errors.Join(err, results[i].Err)
	}

	matchPlayersMap := make(map[string][]entity.Player)
	for _, result := range results {
		if !result.IsNewMatch || result.Err != nil {
			continue
		}

		matchPlayersMap[result.MatchId] = append(matchPlayersMap[result.MatchId], result.Player)
	}

	newMatchesChannel := make(chan error, len(matchPlayersMap))
	for matchId, players := range matchPlayersMap {
		go func(matchId string, players []entity.Player, channel chan error) {
			channel <- u.handleNewMatch(matchId, players)
		}(matchId, players, newMatchesChannel)
	}

	for i := 0; i < len(matchPlayersMap); i++ {
		err = errors.Join(err, <-newMatchesChannel)
	}

	return err
}

func (u PolicePatrol) handleNewMatch(matchId string, players []entity.Player) error {
	matchInfo, err := u.LolService.GetMatchByMatchId(matchId)
	if err != nil {
		return err
	}

	allowedQueueIdTypes := []int{
		int(enum.NormalId),
		int(enum.SoloId),
		int(enum.FlexId),
		int(enum.AramId),
		int(enum.QuickPlayId),
	}
	if !slices.Contains(allowedQueueIdTypes, matchInfo.Info.QueueID) {
		return errors.New("queue id type not supported")
	}

	matchEntity := entity.NewMatch(
		matchId,
		matchInfo.Info.QueueID,
		matchInfo.Info.GameCreation,
		matchInfo.Info.GameEndTimestamp,
		matchInfo.Info.GameDuration,
	)

	err = u.MatchRepository.Create(matchEntity)
	if err == repository.ErrMatchAlreadyExists {
		matchEntity, err = u.MatchRepository.GetByMatchId(matchId)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	var matchParticipantEntities []entity.MatchParticipant
	for _, participant := range matchInfo.Info.Participants {
		for _, player := range players {
			if participant.PUUID == player.Puuid {
				matchParticipantEntity := entity.NewMatchParticipant(
					matchEntity,
					player,
					participant.ChampionName,
					participant.Lane,
					participant.Kills,
					participant.Deaths,
					participant.Assists,
					participant.Win,
				)

				err = errors.Join(err, u.MatchParticipantRepository.Create(matchParticipantEntity))

				matchParticipantEntities = append(matchParticipantEntities, matchParticipantEntity)
			}
		}
	}

	if len(matchParticipantEntities) == 0 {
		return errors.New("cant match player puuid with match participant puuid on match data")
	}

	queueIdTypeMap := map[enum.QueueId]enum.QueueType{
		enum.SoloId: enum.Solo,
		enum.FlexId: enum.Flex,
	}

	matchParticipantEvents := make([]service.MatchParticipantEvent, len(matchParticipantEntities))
	for idx, matchParticipant := range matchParticipantEntities {
		matchParticipantEvent := service.MatchParticipantEvent{
			MatchParticipant: matchParticipant,
		}

		if matchInfo.Info.QueueID == int(enum.SoloId) || matchInfo.Info.QueueID == int(enum.FlexId) {
			leagues, err := u.LolService.GetLeaguesBySummonerId(matchParticipant.Player.SummonerId)
			if err != nil {
				return err
			}
			for _, league := range leagues {
				if enum.QueueType(league.QueueType) == queueIdTypeMap[enum.QueueId(matchInfo.Info.QueueID)] {
					matchParticipantEvent.LeagueItem = &league
					break
				}
			}

			rankedInfo, err := u.RankedInfoRepository.GetByPlayerAndQueueType(
				matchParticipant.Player,
				queueIdTypeMap[enum.QueueId(matchInfo.Info.QueueID)],
			)
			if err == nil {
				matchParticipantEvent.RankedInfo = &rankedInfo
			}

		}

		matchParticipantEvents[idx] = matchParticipantEvent
	}

	templateData := service.NewMatchData{
		Match:                  matchEntity,
		MatchParticipantEvents: matchParticipantEvents,
	}

	message, err := u.TemplateService.ExecuteNewMatchMessageTemplate(templateData)
	if err != nil {
		return err
	}

	whatsappGroupUser := os.Getenv("WPP_GROUP_USER")
	_, err = u.WhatsappService.SendMessageToGroup(message, whatsappGroupUser)
	if err != nil {
		return errors.New("cannot send message to whatsapp group")
	}

	for _, matchParticipantEvent := range matchParticipantEvents {
		if matchParticipantEvent.LeagueItem != nil && matchParticipantEvent.RankedInfo != nil {
			matchParticipantEvent.RankedInfo.LeaguePoints = matchParticipantEvent.LeagueItem.LeaguePoints
			matchParticipantEvent.RankedInfo.Rank = matchParticipantEvent.LeagueItem.Rank
			matchParticipantEvent.RankedInfo.Tier = matchParticipantEvent.LeagueItem.Tier
			matchParticipantEvent.RankedInfo.Wins = matchParticipantEvent.LeagueItem.Wins
			matchParticipantEvent.RankedInfo.Losses = matchParticipantEvent.LeagueItem.Losses

			if _err := u.RankedInfoRepository.Update(*matchParticipantEvent.RankedInfo); _err != nil {
				err = errors.Join(err, _err)
			}
		} else if matchParticipantEvent.LeagueItem != nil && matchParticipantEvent.RankedInfo == nil {
			newRankedInfo := entity.NewRankedInfo(
				matchParticipantEvent.MatchParticipant.Player,
				enum.QueueType(matchParticipantEvent.LeagueItem.QueueType),
				matchParticipantEvent.LeagueItem.Tier,
				matchParticipantEvent.LeagueItem.Rank,
				matchParticipantEvent.LeagueItem.LeaguePoints,
				matchParticipantEvent.LeagueItem.Wins,
				matchParticipantEvent.LeagueItem.Losses,
			)

			if _err := u.RankedInfoRepository.Create(newRankedInfo); _err != nil {
				err = errors.Join(err, _err)
			}
		}
	}

	return err
}

func (u PolicePatrol) checkPlayerHasNewMatch(player entity.Player) (bool, string, error) {
	matchIds, err := u.LolService.GetMatchIdListByPuuid(player.Puuid)
	if err != nil {
		return false, "", err
	}

	if len(matchIds) == 0 {
		return false, "", nil
	}

	remoteLatestMatchId := matchIds[0]

	latestMatch, err := u.MatchRepository.GetLastestByPlayer(player)
	switch err {
	case nil:
		break
	case repository.ErrMatchNotFound:
		return true, remoteLatestMatchId, nil
	default:
		return false, "", err
	}

	if remoteLatestMatchId == latestMatch.MatchId {
		return false, "", nil
	}

	return true, remoteLatestMatchId, nil
}
