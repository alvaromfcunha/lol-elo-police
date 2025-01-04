package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/repository"
	"github.com/alvaromfcunha/lol-elo-police/internal/generated/database"
)

type MatchParticipantData struct {
	Ctx     context.Context
	Queries *database.Queries
}

func NewMatchParticipantData(ctx context.Context, db *sql.DB) MatchParticipantData {
	return MatchParticipantData{
		Ctx:     ctx,
		Queries: database.New(db),
	}
}

func (d MatchParticipantData) Create(mp entity.MatchParticipant) error {
	_, err := d.Queries.CreateMatchParticipant(
		d.Ctx,
		database.CreateMatchParticipantParams{
			ExternalID:        mp.Id.String(),
			PlayerExternalID:  mp.Player.Id.String(),
			MatchesExternalID: mp.Match.Id.String(),
			Champion:          mp.Champion,
			Kills:             int64(mp.Kills),
			Deaths:            int64(mp.Deaths),
			Assists:           int64(mp.Assists),
			IsWin:             mp.IsWin,
		},
	)

	if err != nil {
		fmt.Println("Cannot create MatchParticipant:", err.Error())
		return repository.ErrCannotCreateMatchParticipant
	}

	return nil
}
