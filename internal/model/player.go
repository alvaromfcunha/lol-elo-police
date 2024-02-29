package model

import (
	"github.com/alvaromfcunha/lol-elo-police/pkg/lol"
	"gorm.io/gorm"
)

type Player struct {
	gorm.Model
	SummonerId   string
	GameName     string
	TagLine      string
	Tier         lol.Tier
	Rank         lol.Rank
	LeaguePoints int
	Wins         int
	Losses       int
}
