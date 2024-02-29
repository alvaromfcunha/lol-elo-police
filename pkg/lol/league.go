package lol

type QueueType string

const (
	Flex QueueType = "RANKED_FLEX_SR"
	Solo QueueType = "RANKED_SOLO_5x5"
)

type Tier string

const (
	Bronze      Tier = "BRONZE"
	Silver      Tier = "SILVER"
	Gold        Tier = "GOLD"
	Platinum    Tier = "PLATINUM"
	Emerald     Tier = "EMERALD"
	Diamond     Tier = "DIAMOND"
	Master      Tier = "MASTER"
	Grandmaster Tier = "GRANDMASTER"
	Challenger  Tier = "CHALLENGER"
)

type Rank string

const (
	I   Rank = "I"
	II  Rank = "II"
	III Rank = "III"
	IV  Rank = "IV"
)

type LeagueEntry struct {
	QueueType    QueueType
	Tier         Tier
	Rank         Rank
	LeaguePoints int
	Wins         int
	Losses       int
}

func (l LolApi) GetLeaguesBySummonerId(id string) (lgs []LeagueEntry, err error) {
	url := BR1_LOL_URL +
		"/league/v4/entries/by-summoner/" +
		id

	err = DoRequest(l, "GET", url, &lgs)

	return
}
