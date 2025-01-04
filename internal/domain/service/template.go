package service

import (
	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
)

type MatchParticipantEvent struct {
	MatchParticipant entity.MatchParticipant
	RankedInfo       *entity.RankedInfo
	LeagueItem       *lol.LeagueItem
}

type NewMatchData struct {
	Match                  entity.Match
	MatchParticipantEvents []MatchParticipantEvent
}

type ITemplateService interface {
	ExecuteNewMatchMessageTemplate(data NewMatchData) (string, error)
}
