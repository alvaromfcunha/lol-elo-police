package repository

import (
	"errors"

	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity/enum"
)

var ErrCannotCreateRankedInfo = errors.New("cannot create ranked info")
var ErrCannotUpdateRankedInfo = errors.New("cannot update ranked info")
var ErrCannotGetRankedInfo = errors.New("cannot get ranked info")
var ErrNoRankedInfoFound = errors.New("no ranked info found")

type IRankedInfoRepository interface {
	Create(rankedInfo entity.RankedInfo) error
	GetByPlayerAndQueueType(player entity.Player, queueType enum.QueueType) (entity.RankedInfo, error)
	GetLatestByPlayerAndQueueType(player entity.Player, queueType enum.QueueType) (entity.RankedInfo, error)
}
