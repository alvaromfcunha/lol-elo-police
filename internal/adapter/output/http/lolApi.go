package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/alvaromfcunha/lol-elo-police/internal/domain/service"
)

const (
	AMERICAS_RIOT_URL = "https://americas.api.riotgames.com/riot"
	BR1_LOL_URL       = "https://br1.api.riotgames.com/lol"
)

type LolApi struct {
	ApiKey string
}

func NewLolApi(
	apiKey string,
) LolApi {
	return LolApi{
		ApiKey: apiKey,
	}
}

func (s LolApi) fetch(method string, url string, response any) (err error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return
	}

	req.Header.Set("X-Riot-Token", s.ApiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		err = errors.New("unsuccessful status code: " + fmt.Sprint(res.StatusCode))
		return
	}

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return
	}

	return
}

func (s LolApi) GetAccountByRiotId(riotId service.RiotId) (acc service.RiotAccount, err error) {
	url := AMERICAS_RIOT_URL +
		"/account/v1/accounts/by-riot-id/" +
		riotId.GameName + "/" + riotId.TagLine

	err = s.fetch("GET", url, &acc)

	return
}

func (s LolApi) GetLeaguesBySummonerId(id string) (lgs []service.LeagueEntry, err error) {
	url := BR1_LOL_URL +
		"/league/v4/entries/by-summoner/" +
		id

	err = s.fetch("GET", url, &lgs)

	return
}

func (s LolApi) GetSummonerByPuuid(puuid string) (smm service.Summoner, err error) {
	url := BR1_LOL_URL +
		"/summoner/v4/summoners/by-puuid/" +
		puuid

	err = s.fetch("GET", url, &smm)

	return
}
