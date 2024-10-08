package data

import (
	"context"
	"database/sql"

	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
	"github.com/alvaromfcunha/lol-elo-police/internal/generated/database"
)

type PlayerData struct {
	Ctx     context.Context
	Queries *database.Queries
}

func NewPlayerData(ctx context.Context, db *sql.DB) PlayerData {
	return PlayerData{
		Ctx:     ctx,
		Queries: database.New(db),
	}
}

func (r PlayerData) Create(player entity.Player) error {
	_, err := r.Queries.CreatePlayer(
		r.Ctx,
		database.CreatePlayerParams{
			ExternalID: player.Id.String(),
			SummonerID: player.SummonerId,
			Puuid:      player.Puuid,
			GameName:   player.GameName,
			TagLine:    player.TagLine,
		},
	)

	return err
}

func (r PlayerData) GetAll() (player []entity.Player, err error) {
	records, err := r.Queries.GetPlayers(r.Ctx)
	if err != nil {
		return
	}

	for _, record := range records {
		player = append(player, AssemblePlayer(record.Player))
	}

	return
}
