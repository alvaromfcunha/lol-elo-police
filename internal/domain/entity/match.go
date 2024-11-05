package entity

import (
	"time"

	"github.com/google/uuid"
)

type Match struct {
	Id             uuid.UUID
	MatchId        string
	QueueIdType    int
	GameCreationAt time.Time
	GameEndedAt    time.Time
	GameDuration   time.Duration
}

func NewMatch(
	matchId string,
	queueIdType int,
	gameCreationAt int64,
	gameEndedAt int64,
	gameDuration int,
) Match {
	return Match{
		Id:             uuid.New(),
		MatchId:        matchId,
		QueueIdType:    queueIdType,
		GameCreationAt: time.UnixMilli(gameCreationAt),
		GameEndedAt:    time.UnixMilli(gameEndedAt),
		GameDuration:   time.Duration(gameDuration),
	}
}
