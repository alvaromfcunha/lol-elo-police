package repository

import (
	"errors"

	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
)

var ErrCannotCreatePlayer = errors.New("cannot create player")
var ErrPlayerAlreadyExists = errors.New("player already exists")
var ErrCannotGetPlayer = errors.New("cannot get player")

type IPlayerRepository interface {
	Create(player entity.Player) error
	GetAll() ([]entity.Player, error)
}
