package usecase

import (
	"errors"
	"os"
	"slices"
	"sync"

	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/logger"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity/enum"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/repository"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/service"
	"golang.org/x/sync/errgroup"
)

type PolicePatrolUseCase struct {
	PlayerRepository           repository.IPlayerRepository
	RankedInfoRepository       repository.IRankedInfoRepository
	MatchRepository            repository.IMatchRepository
	MatchParticipantRepository repository.IMatchParticipantRepository
	LolService                 service.ILolService
	WhatsappService            service.IWhatsappService
	TemplateService            service.ITemplateService
}

func NewPolicePatrolUseCase(
	playerRepository repository.IPlayerRepository,
	rankedInfoRepository repository.IRankedInfoRepository,
	matchRepository repository.IMatchRepository,
	matchParticipantRepository repository.IMatchParticipantRepository,
	lolService service.ILolService,
	whatsappService service.IWhatsappService,
	templateService service.ITemplateService,
) PolicePatrolUseCase {
	return PolicePatrolUseCase{
		PlayerRepository:           playerRepository,
		RankedInfoRepository:       rankedInfoRepository,
		MatchRepository:            matchRepository,
		MatchParticipantRepository: matchParticipantRepository,
		LolService:                 lolService,
		WhatsappService:            whatsappService,
		TemplateService:            templateService,
	}
}

func (u PolicePatrolUseCase) Execute() error {
	logger.Debug(u, "Executing police patrol use case")

	players, err := u.PlayerRepository.GetAll()
	if err != nil {
		logger.Error(u, "Error getting players", err)
		return errors.New("Cannot get all players on database")
	}

	matches, err := u.getNewMatches(players)
	if err != nil {
		logger.Error(u, "Error getting new matches", err)
		return errors.New("Cannot remotely get new matches")
	}

	allowedQueueIds := []int{
		int(enum.NormalId),
		int(enum.SoloId),
		int(enum.FlexId),
		int(enum.AramId),
		int(enum.QuickPlayId),
		int(enum.SwiftPlayId),
	}

	filteredMatches := slices.DeleteFunc(matches, func(m lol.Match) bool {
		return !slices.Contains(allowedQueueIds, m.Info.QueueID)
	})

	matchesEntities, matchesParticipantsEntities, errs := u.createMatchesEntities(filteredMatches, players)
	if len(errs) > 0 {
		for _, err := range errs {
			logger.Error(u, "Partial error on creating entities", err)
		}
		if len(errs) == len(filteredMatches) {
			return errors.Join(errs...)
		}
	}

	if errs := u.sendMessages(matchesEntities, matchesParticipantsEntities); len(errs) > 0 {
		for _, err := range errs {
			logger.Error(u, "Partial error on sending messages", err)
		}
		if len(errs) == len(filteredMatches) {
			return errors.Join(errs...)
		}
	}

	return nil
}

func (u PolicePatrolUseCase) getPlayerLastRemoteMatch(player entity.Player) (string, error) {
	matchIds, err := u.LolService.GetMatchIdListByPuuid(player.Puuid)
	if err != nil {
		return "", err
	}

	if len(matchIds) == 0 {
		return "", nil
	}

	remoteLatestMatchId := matchIds[0]

	return remoteLatestMatchId, nil
}

func (u PolicePatrolUseCase) createMatchesEntities(matches []lol.Match, players []entity.Player) ([]entity.Match, [][]entity.MatchParticipant, []error) {
	matchesEntities := []entity.Match{}
	matchesParticipantsEntities := [][]entity.MatchParticipant{}

	wg := sync.WaitGroup{}

	type createMatchEntitiesReturn struct {
		match        entity.Match
		participants []entity.MatchParticipant
	}
	createMatchChan := make(chan createMatchEntitiesReturn, len(matches))
	errorsChan := make(chan error, len(matches))

	for idx := range matches {
		wg.Add(1)
		go func() {
			defer wg.Done()

			matchEntity, participantEntities, err := u.createMatchEntities(matches[idx], players)
			if err != nil {
				errorsChan <- err
				return
			}

			createMatchChan <- createMatchEntitiesReturn{matchEntity, participantEntities}
		}()
	}

	wg.Wait()
	close(createMatchChan)
	close(errorsChan)

	for r := range createMatchChan {
		matchesEntities = append(matchesEntities, r.match)
		matchesParticipantsEntities = append(matchesParticipantsEntities, r.participants)
	}

	errs := []error{}
	for err := range errorsChan {
		errs = append(errs, err)
	}

	return matchesEntities, matchesParticipantsEntities, errs
}

func (u PolicePatrolUseCase) createMatchEntities(match lol.Match, players []entity.Player) (entity.Match, []entity.MatchParticipant, error) {
	participantEntities := []entity.MatchParticipant{}

	matchEntity := entity.NewMatch(
		match.Metadata.MatchID,
		match.Info.QueueID,
		match.Info.GameCreation,
		match.Info.GameEndTimestamp,
		match.Info.GameDuration,
	)

	err := u.MatchRepository.Create(matchEntity)
	if err == repository.ErrMatchAlreadyExists {
		matchEntity, err = u.MatchRepository.GetByMatchId(match.Metadata.MatchID)
		if err != nil {
			return matchEntity, participantEntities, err
		}
	} else if err != nil {
		return matchEntity, participantEntities, err
	}

	participants := slices.DeleteFunc(match.Info.Participants, func(pr *lol.Participant) bool {
		return !slices.ContainsFunc(players, func(pl entity.Player) bool {
			return pl.Puuid == pr.PUUID
		})
	})

	for _, participant := range participants {
		pIdx := slices.IndexFunc(players, func(p entity.Player) bool {
			return p.Puuid == participant.PUUID
		})
		player := players[pIdx]

		var newRankedInfo *entity.RankedInfo
		var prevRankedInfo *entity.RankedInfo
		if matchEntity.QueueIdType == int(enum.SoloId) || matchEntity.QueueIdType == int(enum.FlexId) {
			queueType := enum.QueueIdTypeMap[enum.QueueId(matchEntity.QueueIdType)]

			if ori, err := u.RankedInfoRepository.GetLatestByPlayerAndQueueType(player, queueType); err == nil {
				prevRankedInfo = &ori
			}

			leagueItem, err := u.getPlayerRankedInfo(player, queueType)
			if err != nil {
				return matchEntity, participantEntities, err
			}

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

				if err := u.RankedInfoRepository.Create(nri); err != nil {
					return matchEntity, participantEntities, err
				}

				newRankedInfo = &nri
			}
		}

		participant := entity.NewMatchParticipant(
			matchEntity,
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

		if err := u.MatchParticipantRepository.Create(participant); err != nil {
			return matchEntity, participantEntities, err
		}

		participantEntities = append(participantEntities, participant)
	}

	return matchEntity, participantEntities, nil
}

func (u PolicePatrolUseCase) getPlayerRankedInfo(player entity.Player, queueType enum.QueueType) (*lol.LeagueItem, error) {
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

func (u PolicePatrolUseCase) sendMessages(matches []entity.Match, matchesParticipants [][]entity.MatchParticipant) []error {
	wg := sync.WaitGroup{}
	errChan := make(chan error, len(matches))

	for idx := range matches {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := u.sendMessage(matches[idx], matchesParticipants[idx]); err != nil {
				errChan <- err
			}
		}()
	}

	wg.Wait()
	close(errChan)

	errs := []error{}
	for err := range errChan {
		errs = append(errs, err)
	}

	return errs
}

func (u PolicePatrolUseCase) sendMessage(match entity.Match, participants []entity.MatchParticipant) error {
	filteredParticipants := slices.DeleteFunc(participants, func(p entity.MatchParticipant) bool {
		return !slices.Contains(p.Player.NotifyQueues, enum.QueueId(match.QueueIdType))
	})

	if len(filteredParticipants) == 0 {
		return nil
	}

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

func (u PolicePatrolUseCase) getNewMatches(players []entity.Player) ([]lol.Match, error) {
	matches := []lol.Match{}

	matchIdEg := errgroup.Group{}
	matchIdChan := make(chan string, len(players))

	for _, player := range players {
		matchIdEg.Go(func() error {
			remoteMatchId, err := u.getPlayerLastRemoteMatch(player)
			if err != nil {
				return err
			}

			localMatch, err := u.MatchRepository.GetLastestByPlayer(player)
			if err == repository.ErrMatchNotFound {
				matchIdChan <- remoteMatchId
				return nil
			} else if err != nil {
				return err
			}

			if remoteMatchId != localMatch.MatchId {
				matchIdChan <- remoteMatchId
			}

			return nil
		})
	}

	err := matchIdEg.Wait()
	if err != nil {
		return matches, err
	}

	close(matchIdChan)

	matchEg := errgroup.Group{}
	matchChan := make(chan lol.Match, len(matchIdChan))

	for matchId := range matchIdChan {
		matchEg.Go(func() error {
			match, err := u.LolService.GetMatchByMatchId(matchId)
			if err != nil {
				return err
			}

			matchChan <- match
			return nil
		})
	}

	err = matchEg.Wait()
	if err != nil {
		return matches, err
	}

	close(matchChan)

	for match := range matchChan {
		matches = append(matches, match)
	}

	return matches, nil
}
