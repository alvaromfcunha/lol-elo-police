package entity

import "github.com/google/uuid"

type MatchParticipant struct {
	Id       uuid.UUID
	Match    Match
	Player   Player
	NewRankedInfo *RankedInfo
	PrevRankedInfo *RankedInfo
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
	newRankedInfo *RankedInfo,
	prevRankedInfo *RankedInfo,
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
		NewRankedInfo: newRankedInfo,
		PrevRankedInfo: prevRankedInfo,
		Champion: champion,
		Role:     role,
		Kills:    kills,
		Deaths:   deaths,
		Assists:  assists,
		IsWin:    isWin,
	}
}
