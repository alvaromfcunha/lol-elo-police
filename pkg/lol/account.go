package lol

type RiotId struct {
	GameName string
	TagLine  string
}

type RiotAccount struct {
	Puuid    string
	GameName string
	TagLine  string
}

func (l LolApi) GetAccountByRiotId(riotId RiotId) (acc RiotAccount, err error) {
	url := AMERICAS_RIOT_URL +
		"/account/v1/accounts/by-riot-id/" +
		riotId.GameName + "/" + riotId.TagLine

	err = DoRequest(l, "GET", url, &acc)

	return
}
