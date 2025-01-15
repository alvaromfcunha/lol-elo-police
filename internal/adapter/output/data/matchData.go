package data

import (
	"context"
	"database/sql"

	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/logger"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/repository"
	"github.com/alvaromfcunha/lol-elo-police/internal/generated/database"
	"github.com/mattn/go-sqlite3"
)

type MatchData struct {
	Ctx     context.Context
	Queries *database.Queries
}

func NewMatchData(ctx context.Context, db *sql.DB) MatchData {
	return MatchData{
		Ctx:     ctx,
		Queries: database.New(db),
	}
}

func (r MatchData) Create(match entity.Match) error {
	logger.Debug(r, "Creating match")

	_, err := r.Queries.CreateMatch(
		r.Ctx,
		database.CreateMatchParams{
			ExternalID:     match.Id.String(),
			MatchID:        match.MatchId,
			QueueIDType:    int64(match.QueueIdType),
			GameCreationAt: match.GameCreationAt,
			GameEndedAt:    match.GameEndedAt,
			GameDuration:   match.GameDuration.Milliseconds(),
		},
	)

	switch err := err.(type) {
	case nil:
		break
	case sqlite3.Error:
		if err.Code == 19 {
			return repository.ErrMatchAlreadyExists
		}
	default:
		logger.Error(r, "Cannot create match", err)
		return repository.ErrCannotCreateMatch
	}

	return nil
}

func (r MatchData) GetLastestByPlayer(player entity.Player) (entity.Match, error) {
	logger.Debug(r, "Getting latest match by player")

	record, err := r.Queries.GetLastestMatchesByPlayerExternalId(
		r.Ctx,
		player.Id.String(),
	)

	var match entity.Match
	if err != nil {
		logger.Error(r, "Cannot get latest match by player", err)
		return match, repository.ErrCannotGetMatch
	}
	if len(record) == 0 {
		return match, repository.ErrMatchNotFound
	}

	match = AssembleMatch(record[0].Match)

	return match, nil
}

func (r MatchData) GetByMatchId(matchId string) (entity.Match, error) {
	logger.Debug(r, "Getting match by id")

	record, err := r.Queries.GetMatchesByMatchId(
		r.Ctx,
		matchId,
	)

	var match entity.Match
	if err == sql.ErrNoRows {
		return match, repository.ErrMatchNotFound
	} else if err != nil {
		logger.Error(r, "Cannot get match by match id", err)
		return match, repository.ErrCannotGetMatch
	}

	match = AssembleMatch(record.Match)

	return match, nil
}
