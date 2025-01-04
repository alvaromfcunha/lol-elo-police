package entity

import "github.com/google/uuid"

type MatchParticipant struct {
	Id       uuid.UUID
	Match    Match
	Player   Player
	Champion string
	Role     string
	Kills    int
	Deaths   int
	Assists  int
	IsWin    bool
}

func NewMatchParticipant(
	match Match,
	player Player,
	champion string,
	role string,
	kills int,
	deaths int,
	assists int,
	isWin bool,
) MatchParticipant {
	return MatchParticipant{
		Id:       uuid.New(),
		Match:    match,
		Player:   player,
		Champion: champion,
		Role:     role,
		Kills:    kills,
		Deaths:   deaths,
		Assists:  assists,
		IsWin:    isWin,
	}
}
