package service

import (
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
)

type ITemplateService interface {
	ExecuteNewMatchMessageTemplate(match entity.Match, participants []entity.MatchParticipant) (string, error)
}
