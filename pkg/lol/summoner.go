package lol

type Summoner struct {
	Id string
}

func (l LolApi) GetSummonerByPuuid(puuid string) (smm Summoner, err error) {
	url := BR1_LOL_URL +
		"/summoner/v4/summoners/by-puuid/" +
		puuid

	err = DoRequest(l, "GET", url, &smm)

	return
}
