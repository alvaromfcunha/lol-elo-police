package repository

import "github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"

type IMatchRepository interface {
	Create(match entity.Match) error
	Update(match entity.Match) error
	GetLastestByPlayer(player entity.Player) (entity.Match, error)
}
