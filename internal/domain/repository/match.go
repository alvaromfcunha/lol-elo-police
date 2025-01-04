package repository

import (
	"errors"

	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
)

var ErrCannotCreateMatch = errors.New("cannot create match")
var ErrMatchAlreadyExists = errors.New("match already exists")
var ErrCannotGetMatch = errors.New("cannot get match")
var ErrMatchNotFound = errors.New("match not found")

type IMatchRepository interface {
	Create(match entity.Match) error
	GetLastestByPlayer(player entity.Player) (entity.Match, error)
	GetByMatchId(matchId string) (entity.Match, error)
}
