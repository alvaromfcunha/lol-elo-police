package repository

import (
	"errors"

	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
)

var ErrCannotCreateMatchParticipant = errors.New("cannot create match participant")

// var ErrCannotUpdateMatchParticipant = errors.	New("cannot update match participant")

type IMatchParticipantRepository interface {
	Create(matchParticipant entity.MatchParticipant) error
	// Update(matchParticipant entity.MatchParticipant) error
}
