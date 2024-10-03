package entity

import "github.com/google/uuid"

type Player struct {
	Id         uuid.UUID `json:"id"`
	SummonerId string    `json:"summonerId"`
	GameName   string    `json:"gameName"`
	TagLine    string    `json:"tagLine"`
}

func NewPlayer(
	summonerId string,
	gameName string,
	tagLine string,
) Player {
	return Player{
		Id:         uuid.New(),
		SummonerId: summonerId,
		GameName:   gameName,
		TagLine:    tagLine,
	}
}
