package entity

import (
	"time"

	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity/enum"
	"github.com/google/uuid"
)

type RankedInfo struct {
	Id           uuid.UUID      `json:"id"`
	Player       Player         `json:"-"`
	QueueType    enum.QueueType `json:"queueType"`
	Tier         string         `json:"tier"`
	Rank         string         `json:"rank"`
	LeaguePoints int            `json:"leaguePoints"`
	Wins         int            `json:"wins"`
	Losses       int            `json:"losses"`
	CreatedAt    time.Time      `json:"createdAt"`
}

func NewRankedInfo(
	player Player,
	queueType enum.QueueType,
	tier string,
	rank string,
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
		CreatedAt:    time.Now(),
	}
}
