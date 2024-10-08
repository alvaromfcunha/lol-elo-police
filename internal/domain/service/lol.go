package service

import (
	"github.com/KnutZuidema/golio/riot/account"
	"github.com/KnutZuidema/golio/riot/lol"
)

type ILolService interface {
	GetAccountByRiotId(gameName string, tagLine string) (*account.Account, error)
	GetLeaguesBySummonerId(id string) ([]*lol.LeagueItem, error)
	GetSummonerByPuuid(puuid string) (*lol.Summoner, error)
	GetMatchIdListByPuuid(puuid string) ([]*string, error)
	GetMatchByMatchId(matchId string) (*lol.Match, error)
}
