package service

import (
	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
)

type QueueUpdateData struct {
	Player         entity.Player
	RankedType     string
	NewLeagueEntry lol.LeagueItem
	OldRankedInfo  entity.RankedInfo
}

type QueueNewEntryData struct {
	Player      entity.Player
	RankedType  string
	LeagueEntry lol.LeagueItem
}

type ITemplateService interface {
	ExecuteQueueUpdateMessageTemplate(data QueueUpdateData) (string, error)
	ExecuteQueueNewEntryMessageTemplate(data QueueNewEntryData) (string, error)
}
