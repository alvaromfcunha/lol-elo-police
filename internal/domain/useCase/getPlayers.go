package usecase

import (
	"errors"

	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/logger"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity/enum"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/repository"
)

type GetPlayersUseCase struct {
	PlayerRepository     repository.IPlayerRepository
	RankedInfoRepository repository.IRankedInfoRepository
}

func NewGetPlayersUseCase(
	playerRepository repository.IPlayerRepository,
	rankedInfoRepository repository.IRankedInfoRepository,
) GetPlayersUseCase {
	return GetPlayersUseCase{
		PlayerRepository:     playerRepository,
		RankedInfoRepository: rankedInfoRepository,
	}
}

type playerFacade struct {
	entity.Player

	SoloQueueInfo *entity.RankedInfo `json:"soloQueue"`
	FlexQueueInfo *entity.RankedInfo `json:"flexQueue"`
}

type GetPlayersOutput []playerFacade

func (u GetPlayersUseCase) Execute() (output GetPlayersOutput, err error) {
	logger.Debug(u, "Executing get players use case")

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
