package usecase

import (
	"errors"
	"fmt"
	"os"

	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity/enum"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/repository"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/service"
)

type CreatePlayer struct {
	PlayerRepository     repository.IPlayerRepository
	RankedInfoRepository repository.IRankedInfoRepository
	LolService           service.ILolService
}

type CreatePlayerInput struct {
	GameName string
	TagLine  string
}

type CreatePlayerOutput struct {
	entity.Player
}

func (u CreatePlayer) Execute(input CreatePlayerInput) (output CreatePlayerOutput, err error) {
	account, err := u.LolService.GetAccountByRiotId(
		service.RiotId{
			GameName: input.GameName,
			TagLine:  input.TagLine,
		},
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot get riot account by riot id:", err)
		err = errors.New("cannot get riot account by riot id")
		return
	}

	summoner, err := u.LolService.GetSummonerByPuuid(account.Puuid)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot get summoner by puuid:", err)
		err = errors.New("cannot get summoner by puuid")
		return
	}

	player := entity.NewPlayer(
		summoner.Id,
		account.GameName,
		account.TagLine,
	)

	leagues, err := u.LolService.GetLeaguesBySummonerId(summoner.Id)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot get leagues by summoner id:", err)
		err = errors.New("cannot get leagues by summoner id")
		return
	}

	var soloQueue *service.LeagueEntry
	var flexQueue *service.LeagueEntry
	for _, league := range leagues {
		switch league.QueueType {
		case enum.Solo:
			soloQueue = &league
		case enum.Flex:
			flexQueue = &league
		}
	}

	err = u.PlayerRepository.Create(player)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot create player on database:", err)
		err = errors.New("cannot create player on database")
		return
	}

	if soloQueue != nil {
		soloQueueInfo := entity.NewRankedInfo(
			player,
			enum.Solo,
			soloQueue.Tier,
			soloQueue.Rank,
			soloQueue.LeaguePoints,
			soloQueue.Wins,
			soloQueue.Losses,
		)

		err = u.RankedInfoRepository.Create(soloQueueInfo, player)
		if err != nil {
			fmt.Fprintln(os.Stderr, "cannot solo queue info on database:", err)
			err = errors.New("cannot solo queue info on database")
			return
		}
	}
	if flexQueue != nil {
		flexQueueInfo := entity.NewRankedInfo(
			player,
			enum.Flex,
			flexQueue.Tier,
			flexQueue.Rank,
			flexQueue.LeaguePoints,
			flexQueue.Wins,
			flexQueue.Losses,
		)

		err = u.RankedInfoRepository.Create(flexQueueInfo, player)
		if err != nil {
			fmt.Fprintln(os.Stderr, "cannot flex queue info on database:", err)
			err = errors.New("cannot flex queue info on database")
			return
		}
	}

	output = CreatePlayerOutput{player}

	return
}