package usecase

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity/enum"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/repository"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/service"
)

type PolicePatrolV2 struct {
	PlayerRepository           repository.IPlayerRepository
	RankedInfoRepository       repository.IRankedInfoRepository
	MatchRepository            repository.IMatchRepository
	MatchParticipantRepository repository.IMatchParticipantRepository
	LolService                 service.ILolService
	WhatsappService            service.IWhatsappService
	TemplateService            service.ITemplateService
}

func (u PolicePatrolV2) Execute() error {
	players, err := u.PlayerRepository.GetAll()
	if err != nil {
		return errors.New("cannot get all players on database")
	}

	type checkPlayerHasNewMatch struct {
		IsNewMatch bool
		MatchId    string
		Err        error
	}

	checkPlayersChannel := make([]chan checkPlayerHasNewMatch, len(players))
	for i, player := range players {
		go func(player entity.Player, channel chan checkPlayerHasNewMatch) {
			isNewMatch, matchId, err := u.checkPlayerHasNewMatch(player)
			channel <- checkPlayerHasNewMatch{isNewMatch, matchId, err}
		}(player, checkPlayersChannel[i])
	}

	results := make([]checkPlayerHasNewMatch, len(players))
	for i := range checkPlayersChannel {
		results[i] = <-checkPlayersChannel[i]
		err = errors.Join(err, results[i].Err)
	}

	matchPlayersMap := make(map[string][]entity.Player)
	for i, result := range results {
		if !result.IsNewMatch || result.Err != nil {
			continue
		}

		matchPlayersMap[result.MatchId] = append(matchPlayersMap[result.MatchId], players[i])
	}

	i := 0
	newMatchesChannel := make([]chan error, len(matchPlayersMap))
	for matchId, players := range matchPlayersMap {
		go func(matchId string, players []entity.Player, channel chan error) {
			channel <- u.handleNewMatch(matchId, players)
		}(matchId, players, newMatchesChannel[i])
		i++
	}

	for _, channel := range newMatchesChannel {
		err = errors.Join(err, <-channel)
	}

	return err
}

func (u PolicePatrolV2) handleNewMatch(matchId string, players []entity.Player) error {
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
	if err != nil {
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
					participant.Role,
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

	var message string
	if matchInfo.Info.QueueID == int(enum.SoloId) || matchInfo.Info.QueueID == int(enum.FlexId) {
		matchParticipantsWithLeagueItem := make([]service.MatchParticipantWithRankedInfo, len(matchParticipantEntities))
		for _, matchParticipant := range matchParticipantEntities {
			leagues, err := u.LolService.GetLeaguesBySummonerId(matchParticipant.Player.SummonerId)
			if err != nil {
				return err
			}

			rankedInfo, err := u.RankedInfoRepository.GetByPlayerAndQueueType(
				matchParticipant.Player,
				enum.QueueType(matchInfo.Info.QueueID),
			)
			if err != nil {
				return err
			}

			for _, league := range leagues {
				if league.QueueType == fmt.Sprint(matchInfo.Info.QueueID) {
					matchParticipantWithLeagueItem := service.MatchParticipantWithRankedInfo{
						MatchParticipant: matchParticipant,
						RankedInfo:       rankedInfo,
						LeagueItem:       league,
					}

					matchParticipantsWithLeagueItem = append(matchParticipantsWithLeagueItem, matchParticipantWithLeagueItem)
				}
			}
		}

		templateData := service.NewRankedMatchData{
			Match:                           matchEntity,
			MatchParticipantsWithLeagueItem: matchParticipantsWithLeagueItem,
		}

		message, err = u.TemplateService.ExecuteNewRankedMatchMessageTemplate(templateData)
		if err != nil {
			return err
		}
	} else {
		templateData := service.NewUnrankedMatchData{
			Match:             matchEntity,
			MatchParticipants: matchParticipantEntities,
		}

		message, err = u.TemplateService.ExecuteNewUnrankedMatchMessageTemplate(templateData)
		if err != nil {
			return err
		}
	}

	whatsappGroupUser := os.Getenv("WPP_GROUP_USER")
	_, err = u.WhatsappService.SendMessageToGroup(message, whatsappGroupUser)
	if err != nil {
		return errors.New("cannot send message to whatsapp group")
	}

	return err
}

func (u PolicePatrolV2) checkPlayerHasNewMatch(player entity.Player) (bool, string, error) {
	matchIds, err := u.LolService.GetMatchIdListByPuuid(player.Puuid)
	if err != nil {
		return false, "", err
	}

	latestMatch, err := u.MatchRepository.GetLastestByPlayer(player)
	if err != nil {
		return false, "", err
	}

	remoteLatestMatchId := *matchIds[len(matchIds)-1]

	if remoteLatestMatchId == latestMatch.MatchId {
		return false, "", nil
	}

	return true, remoteLatestMatchId, nil
}
