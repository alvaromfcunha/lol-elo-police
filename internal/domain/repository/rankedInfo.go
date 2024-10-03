package repository

import (
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity/enum"
)

type IRankedInfoRepository interface {
	Create(rankedInfo entity.RankedInfo, player entity.Player) error
	Update(rankedInfo entity.RankedInfo) error
	GetByPlayerAndQueueType(player entity.Player, queueType enum.QueueType) (entity.RankedInfo, error)
}
