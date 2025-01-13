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
	params := database.CreateMatchParticipantParams{
		ExternalID:        mp.Id.String(),
		PlayerExternalID:  mp.Player.Id.String(),
		MatchesExternalID: mp.Match.Id.String(),
		Champion:          mp.Champion,
		Kills:             int64(mp.Kills),
		Deaths:            int64(mp.Deaths),
		Assists:           int64(mp.Assists),
		IsWin:             mp.IsWin,
	}

	if mp.NewRankedInfo != nil {
		params.NewRankedInfoExternalID = mp.NewRankedInfo.Id.String()
	}
	if mp.PrevRankedInfo != nil {
		params.PrevRankedInfoExternalID = mp.PrevRankedInfo.Id.String()
	}

	_, err := d.Queries.CreateMatchParticipant(
		d.Ctx,
		params,
	)

	if err != nil {
		fmt.Println("Cannot create MatchParticipant:", err.Error())
		return repository.ErrCannotCreateMatchParticipant
	}

	return nil
}
