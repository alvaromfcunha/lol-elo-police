package usecase

import (
	"errors"
	"os"
	"slices"

	"github.com/KnutZuidema/golio/riot/lol"
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

	var errs error

	type getPlayerLastRemoteMatchReturn struct {
		Player entity.Player
		MatchId string
		Err     error
	}
	matchIdChan := make(chan getPlayerLastRemoteMatchReturn, len(players))
	for _, player := range players {
		go func(player entity.Player, channel chan getPlayerLastRemoteMatchReturn) {
			matchId, err := u.getPlayerLastRemoteMatch(player)
			channel <- getPlayerLastRemoteMatchReturn{player, matchId, err}
		}(player, matchIdChan)
	}

	matchPlayersMap := map[string][]entity.Player{}
	for range players {
		r := <-matchIdChan
		if r.MatchId != "" {
			matchPlayersMap[r.MatchId] = append(matchPlayersMap[r.MatchId], r.Player)
		}
		errs = errors.Join(errs, r.Err)
	}

	type createEntitiesReturn struct {
		Match             entity.Match
		MatchParticipants []entity.MatchParticipant
		Err             error
	}
	entitiesChan := make(chan createEntitiesReturn, len(matchPlayersMap))
	for matchId, players := range matchPlayersMap {
		go func(channel chan createEntitiesReturn, matchId string, players []entity.Player) {
			match, participants, err := u.createEntries(matchId, players)
			channel <- createEntitiesReturn{match, participants, err}
		}(entitiesChan, matchId, players)
	}

	messagesChan := make(chan error, len(matchPlayersMap))
	for range matchPlayersMap {
		r := <-entitiesChan
		errs = errors.Join(errs, r.Err)
		go func(channel chan error, match entity.Match, participants []entity.MatchParticipant) {
			channel <- u.sendMessage(match, participants)
		}(messagesChan, r.Match, r.MatchParticipants)
	}

	for range matchPlayersMap {
		errors.Join(errs, <-messagesChan)
	}

	return errs
}

func (u PolicePatrol) getPlayerLastRemoteMatch(player entity.Player) (string, error) {
	matchIds, err := u.LolService.GetMatchIdListByPuuid(player.Puuid)
	if err != nil {
		return "", err
	}

	if len(matchIds) == 0 {
		return "", nil
	}

	remoteLatestMatchId := matchIds[0]

	latestMatch, err := u.MatchRepository.GetLastestByPlayer(player)
	switch err {
	case nil:
		break
	case repository.ErrMatchNotFound:
		return remoteLatestMatchId, nil
	default:
		return "", err
	}

	if remoteLatestMatchId == latestMatch.MatchId {
		return "", nil
	}

	return remoteLatestMatchId, nil
}

func (u PolicePatrol) createEntries(matchId string, players []entity.Player) (entity.Match, []entity.MatchParticipant, error) {
	var match entity.Match
	var participants []entity.MatchParticipant

	matchInfo, err := u.LolService.GetMatchByMatchId(matchId)
	if err != nil {
		return match, participants, err
	}

	allowedQueueIdTypes := []int{
		int(enum.NormalId),
		int(enum.SoloId),
		int(enum.FlexId),
		int(enum.AramId),
		int(enum.QuickPlayId),
	}
	if !slices.Contains(allowedQueueIdTypes, matchInfo.Info.QueueID) {
		return match, participants, errors.New("queue id type not supported")
	}

	match = entity.NewMatch(
		matchId,
		matchInfo.Info.QueueID,
		matchInfo.Info.GameCreation,
		matchInfo.Info.GameEndTimestamp,
		matchInfo.Info.GameDuration,
	)

	err = u.MatchRepository.Create(match)
	if err == repository.ErrMatchAlreadyExists {
		match, err = u.MatchRepository.GetByMatchId(matchId)
		if err != nil {
			return match, participants, err
		}
	} else if err != nil {
		return match, participants, err
	}

	var errs error
	for _, participant := range matchInfo.Info.Participants {
		for _, player := range players {
			if participant.PUUID == player.Puuid {

				var newRankedInfo *entity.RankedInfo
				var prevRankedInfo *entity.RankedInfo
				if match.QueueIdType == int(enum.SoloId) || match.QueueIdType == int(enum.FlexId) {
					queueType := enum.QueueIdTypeMap[enum.QueueId(match.QueueIdType)]

					if ori, err := u.RankedInfoRepository.GetLatestByPlayerAndQueueType(player, queueType); err == nil {
						prevRankedInfo = &ori
					}

					leagueItem, err := u.getPlayerRankedInfo(player, queueType)
					errs = errors.Join(errs, err)

					if leagueItem != nil {
						nri := entity.NewRankedInfo(
							player,
							enum.QueueType(leagueItem.QueueType),
							leagueItem.Tier,
							leagueItem.Rank,
							leagueItem.LeaguePoints,
							leagueItem.Wins,
							leagueItem.Losses,
						)

						if err := u.RankedInfoRepository.Create(nri); err == nil {
							newRankedInfo = &nri
						} else {
							errs = errors.Join(errs, err)
						}
					}
				}

				participant := entity.NewMatchParticipant(
					match,
					player,
					newRankedInfo,
					prevRankedInfo,
					participant.ChampionName,
					participant.Lane,
					participant.Kills,
					participant.Deaths,
					participant.Assists,
					participant.Win,
				)

				errs = errors.Join(errs, u.MatchParticipantRepository.Create(participant))

				participants = append(participants, participant)
			}
		}
	}

	return match, participants, errs
}

func (u PolicePatrol) getPlayerRankedInfo(player entity.Player, queueType enum.QueueType) (*lol.LeagueItem, error) {
	lis, err := u.LolService.GetLeaguesBySummonerId(player.SummonerId)
	if err != nil {
		return nil, err
	}

	idx := slices.IndexFunc(lis, func(li lol.LeagueItem) bool {
		return li.QueueType == string(queueType)
	})
	if idx == -1 {
		return nil, nil
	}

	return &lis[idx], nil
}

func (u PolicePatrol) sendMessage(match entity.Match, participants []entity.MatchParticipant) error {
	filteredParticipants := slices.DeleteFunc(participants, func (p entity.MatchParticipant) bool {
		return !slices.Contains(p.Player.NotifyQueues, enum.QueueId(match.QueueIdType))
	})

	message, err := u.TemplateService.ExecuteNewMatchMessageTemplate(match, filteredParticipants)
	if err != nil {
		return err
	}

	whatsappGroupUser := os.Getenv("WPP_GROUP_USER")
	_, err = u.WhatsappService.SendMessageToGroup(message, whatsappGroupUser)
	if err != nil {
		return errors.New("cannot send message to whatsapp group")
	}

	return nil
}
