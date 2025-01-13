package usecase

import (
	"errors"
	"slices"

	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity/enum"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/repository"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/service"
)

type CreatePlayer struct {
	PlayerRepository           repository.IPlayerRepository
	RankedInfoRepository       repository.IRankedInfoRepository
	MatchRepository            repository.IMatchRepository
	MatchParticipantRepository repository.IMatchParticipantRepository
	LolService                 service.ILolService
}

type CreatePlayerInput struct {
	GameName     string
	TagLine      string
	NotifyQueues []enum.QueueId
}

type CreatePlayerOutput struct {
	entity.Player
}

func (u CreatePlayer) Execute(input CreatePlayerInput) (CreatePlayerOutput, error) {
	var output CreatePlayerOutput

	account, err := u.LolService.GetAccountByRiotId(
		input.GameName,
		input.TagLine,
	)
	if err != nil {
		return output, errors.New("cannot get riot account by riot id")
	}

	summoner, err := u.LolService.GetSummonerByPuuid(account.Puuid)
	if err != nil {
		return output, errors.New("cannot get summoner by puuid")
	}

	playerEntity := entity.NewPlayer(
		summoner.ID,
		account.Puuid,
		account.GameName,
		account.TagLine,
		input.NotifyQueues,
	)

	err = u.PlayerRepository.Create(playerEntity)
	if err != nil {
		return output, errors.New("cannot create player on database")
	}

	leagues, err := u.LolService.GetLeaguesBySummonerId(summoner.ID)
	if err != nil {
		return output, errors.New("cannot get leagues by summoner id")
	}

	var soloQueue *lol.LeagueItem
	var flexQueue *lol.LeagueItem
	for _, league := range leagues {
		switch league.QueueType {
		case string(lol.QueueRankedSolo):
			soloQueue = &league
		case string(lol.QueueRankedFlex):
			flexQueue = &league
		}
	}

	var soloQueueEntity *entity.RankedInfo
	var flexQueueEntity *entity.RankedInfo
	if soloQueue != nil {
		sqi := entity.NewRankedInfo(
			playerEntity,
			enum.Solo,
			soloQueue.Tier,
			soloQueue.Rank,
			soloQueue.LeaguePoints,
			soloQueue.Wins,
			soloQueue.Losses,
		)

		err = u.RankedInfoRepository.Create(sqi)
		if err != nil {
			return output, errors.New("cannot solo queue info on database")
		}

		soloQueueEntity = &sqi
	}
	if flexQueue != nil {
		fqi := entity.NewRankedInfo(
			playerEntity,
			enum.Flex,
			flexQueue.Tier,
			flexQueue.Rank,
			flexQueue.LeaguePoints,
			flexQueue.Wins,
			flexQueue.Losses,
		)

		err = u.RankedInfoRepository.Create(fqi)
		if err != nil {
			return output, errors.New("cannot flex queue info on database")
		}

		flexQueueEntity = &fqi
	}

	matchIds, err := u.LolService.GetMatchIdListByPuuid(playerEntity.Puuid)
	if err != nil {
		return output, errors.New("cannot get player matchId list")
	}

	if len(matchIds) != 0 {
		lastMatchId := matchIds[0]
		if lastMatch, err := u.LolService.GetMatchByMatchId(lastMatchId); err == nil {
			matchEntity := entity.NewMatch(
				lastMatch.Metadata.MatchID,
				lastMatch.Info.QueueID,
				lastMatch.Info.GameCreation,
				lastMatch.Info.GameEndTimestamp,
				lastMatch.Info.GameDuration,
			)

			err = u.MatchRepository.Create(matchEntity)
			if err != nil {
				return output, errors.New("cannot create player last match")
			}

			pIdx := slices.IndexFunc(lastMatch.Info.Participants, func(p *lol.Participant) bool {
				return p.PUUID == playerEntity.Puuid
			})

			if pIdx == -1 {
				return output, errors.New("player match does not contain player participant")
			}

			var matchRankedEntity *entity.RankedInfo
			if lastMatch.Info.QueueID == int(enum.SoloId) {
				matchRankedEntity = soloQueueEntity
			} else if lastMatch.Info.QueueID == int(enum.FlexId) {
				matchRankedEntity = flexQueueEntity
			}

			participantEntity := entity.NewMatchParticipant(
				matchEntity,
				playerEntity,
				matchRankedEntity,
				nil,
				lastMatch.Info.Participants[pIdx].ChampionName,
				lastMatch.Info.Participants[pIdx].Role,
				lastMatch.Info.Participants[pIdx].Kills,
				lastMatch.Info.Participants[pIdx].Deaths,
				lastMatch.Info.Participants[pIdx].Assists,
				lastMatch.Info.Participants[pIdx].Win,
			)

			err = u.MatchParticipantRepository.Create(participantEntity)
			if err != nil {
				return output, errors.New("cannot create player match participant")
			}
		}
	}

	output = CreatePlayerOutput{playerEntity}
	return output, nil
}
