package service

import "github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"

type QueueUpdateData struct {
	Player         entity.Player
	RankedType     string
	NewLeagueEntry LeagueEntry
	OldRankedInfo  entity.RankedInfo
}

type QueueNewEntryData struct {
	Player      entity.Player
	RankedType  string
	LeagueEntry LeagueEntry
}

type ITemplateService interface {
	ExecuteQueueUpdateMessageTemplate(data QueueUpdateData) (string, error)
	ExecuteQueueNewEntryMessageTemplate(data QueueNewEntryData) (string, error)
}
