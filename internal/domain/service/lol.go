package service

import "github.com/alvaromfcunha/lol-elo-police/internal/domain/entity/enum"

type RiotId struct {
	GameName string
	TagLine  string
}

type RiotAccount struct {
	Puuid    string
	GameName string
	TagLine  string
}

type LeagueEntry struct {
	QueueType    enum.QueueType
	Tier         enum.Tier
	Rank         enum.Rank
	LeaguePoints int
	Wins         int
	Losses       int
}

type Summoner struct {
	Id string
}

type ILolService interface {
	GetAccountByRiotId(riotId RiotId) (RiotAccount, error)
	GetLeaguesBySummonerId(id string) ([]LeagueEntry, error)
	GetSummonerByPuuid(puuid string) (Summoner, error)
}
