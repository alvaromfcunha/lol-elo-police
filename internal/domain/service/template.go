package service

import (
	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
)

type MatchParticipantWithRankedInfo struct {
	MatchParticipant entity.MatchParticipant
	RankedInfo       entity.RankedInfo
	LeagueItem       *lol.LeagueItem
}

type NewRankedMatchData struct {
	Match                           entity.Match
	MatchParticipantsWithLeagueItem []MatchParticipantWithRankedInfo
}

type NewUnrankedMatchData struct {
	Match             entity.Match
	MatchParticipants []entity.MatchParticipant
}

type ITemplateService interface {
	ExecuteNewRankedMatchMessageTemplate(data NewRankedMatchData) (string, error)
	ExecuteNewUnrankedMatchMessageTemplate(data NewUnrankedMatchData) (string, error)
}
