package entity

import (
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity/enum"
	"github.com/google/uuid"
)

type RankedInfo struct {
	Id           uuid.UUID      `json:"id"`
	Player       Player         `json:"-"`
	QueueType    enum.QueueType `json:"queueType"`
	Tier         enum.Tier      `json:"tier"`
	Rank         enum.Rank      `json:"rank"`
	LeaguePoints int            `json:"leaguePoints"`
	Wins         int            `json:"wins"`
	Losses       int            `json:"losses"`
}

func NewRankedInfo(
	player Player,
	queueType enum.QueueType,
	tier enum.Tier,
	rank enum.Rank,
	leaguePoints int,
	wins int,
	losses int,
) RankedInfo {
	return RankedInfo{
		Id:           uuid.New(),
		Player:       player,
		QueueType:    queueType,
		Tier:         tier,
		Rank:         rank,
		LeaguePoints: leaguePoints,
		Wins:         wins,
		Losses:       losses,
	}
}
