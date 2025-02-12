package data

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/logger"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/repository"
	"github.com/alvaromfcunha/lol-elo-police/internal/generated/database"
	"github.com/mattn/go-sqlite3"
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
	logger.Debug(r, "Creating player")

	nqs := make([]string, len(player.NotifyQueues))
	for idx, q := range player.NotifyQueues {
		nqs[idx] = strconv.Itoa(int(q))
	}

	_, err := r.Queries.CreatePlayer(
		r.Ctx,
		database.CreatePlayerParams{
			ExternalID: player.Id.String(),
			SummonerID: player.SummonerId,
			Puuid:      player.Puuid,
			GameName:   player.GameName,
			TagLine:    player.TagLine,
			NotifyQueues: strings.Join(nqs, ","),
		},
	)

	switch err := err.(type) {
	case nil:
		break
	case sqlite3.Error:
		if err.Code == 19 {
			return repository.ErrPlayerAlreadyExists
		}
	default:
		logger.Error(r, "Cannot create match", err)
		return repository.ErrCannotCreatePlayer
	}

	return nil
}

func (r PlayerData) GetAll() ([]entity.Player, error) {
	records, err := r.Queries.GetPlayers(r.Ctx)

	var players []entity.Player
	if err != nil {
		return players, repository.ErrCannotGetPlayer
	}

	for _, record := range records {
		players = append(players, AssemblePlayer(record.Player))
	}

	return players, nil
}
