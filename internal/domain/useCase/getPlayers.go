package usecase

import (
	"errors"

	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity/enum"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/repository"
)

type GetPlayers struct {
	PlayerRepository     repository.IPlayerRepository
	RankedInfoRepository repository.IRankedInfoRepository
}

type playerFacade struct {
	entity.Player
	SoloQueueInfo *entity.RankedInfo `json:"soloQueue"`
	FlexQueueInfo *entity.RankedInfo `json:"flexQueue"`
}

type GetPlayersOutput []playerFacade

func (u GetPlayers) Execute() (output GetPlayersOutput, err error) {
	players, err := u.PlayerRepository.GetAll()
	if err != nil {
		err = errors.New("cannot get all players on database")
	}

	output = *new(GetPlayersOutput)
	for _, player := range players {
		var soloQueueInfo *entity.RankedInfo
		if rankedInfo, err := u.RankedInfoRepository.GetByPlayerAndQueueType(player, enum.Solo); err == nil {
			soloQueueInfo = &rankedInfo
		}

		var flexQueueInfo *entity.RankedInfo
		if rankedInfo, err := u.RankedInfoRepository.GetByPlayerAndQueueType(player, enum.Flex); err == nil {
			flexQueueInfo = &rankedInfo
		}

		playerFacade := playerFacade{
			Player:        player,
			SoloQueueInfo: soloQueueInfo,
			FlexQueueInfo: flexQueueInfo,
		}

		output = append(output, playerFacade)
	}

	return
}
