package repository

import "github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"

type IMatchParticipantRepository interface {
	Create(matchParticipant entity.MatchParticipant) error
	Update(matchParticipant entity.MatchParticipant) error
}
