package data

import (
	"context"
	"database/sql"

	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity/enum"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/repository"
	"github.com/alvaromfcunha/lol-elo-police/internal/generated/database"
)

type RankedInfoData struct {
	Ctx     context.Context
	Queries *database.Queries
}

func NewRankedInfoData(ctx context.Context, db *sql.DB) RankedInfoData {
	return RankedInfoData{
		Ctx:     ctx,
		Queries: database.New(db),
	}
}

func (r RankedInfoData) Create(rankedInfo entity.RankedInfo) error {
	_, err := r.Queries.CreateRankedInfo(
		r.Ctx,
		database.CreateRankedInfoParams{
			ExternalID:       rankedInfo.Id.String(),
			PlayerExternalID: rankedInfo.Player.Id.String(),
			QueueType:        string(rankedInfo.QueueType),
			Tier:             string(rankedInfo.Tier),
			Rank:             string(rankedInfo.Rank),
			LeaguePoints:     int64(rankedInfo.LeaguePoints),
			Wins:             int64(rankedInfo.Wins),
			Losses:           int64(rankedInfo.Losses),
		},
	)

	if err != nil {
		return repository.ErrCannotCreateRankedInfo
	}

	return nil
}

func (r RankedInfoData) Update(rankedInfo entity.RankedInfo) error {
	err := r.Queries.UpdateRankedInfo(
		r.Ctx,
		database.UpdateRankedInfoParams{
			ExternalID:   rankedInfo.Id.String(),
			Tier:         string(rankedInfo.Tier),
			Rank:         string(rankedInfo.Rank),
			LeaguePoints: int64(rankedInfo.LeaguePoints),
			Wins:         int64(rankedInfo.Wins),
			Losses:       int64(rankedInfo.Losses),
		},
	)

	if err != nil {
		return repository.ErrCannotUpdateRankedInfo
	}

	return nil
}

func (r RankedInfoData) GetByPlayerAndQueueType(player entity.Player, queueType enum.QueueType) (entity.RankedInfo, error) {
	var rankedInfo entity.RankedInfo

	record, err := r.Queries.GetByPlayerExternalIdAndQueueType(
		r.Ctx,
		database.GetByPlayerExternalIdAndQueueTypeParams{
			PlayerExternalID: player.Id.String(),
			QueueType:        string(queueType),
		},
	)

	if err == sql.ErrNoRows {
		return rankedInfo, repository.ErrNoRankedInfoFound
	} else if err != nil {
		return rankedInfo, repository.ErrCannotGetRankedInfo
	}

	rankedInfo = AssembleRankedInfo(record.RankedInfo, record.Player)

	return rankedInfo, nil
}
