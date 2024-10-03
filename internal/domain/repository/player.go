package repository

import "github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"

type IPlayerRepository interface {
	Create(player entity.Player) error
	GetAll() ([]entity.Player, error)
}
