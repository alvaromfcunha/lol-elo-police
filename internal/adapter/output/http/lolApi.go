package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/KnutZuidema/golio/riot/account"
	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/logger"
)

const (
	AMERICAS_URL = "https://americas.api.riotgames.com"
	BR1_LOL_URL  = "https://br1.api.riotgames.com/lol"
)

type LolApi struct {
	client *RateLimitedClient
	apiKey string
}

func NewLolApi(
	client *RateLimitedClient,
	apiKey string,
) LolApi {
	return LolApi{
		client: client,
		apiKey: apiKey,
	}
}

func (s LolApi) fetch(method string, url string, response any) (err error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return
	}

	req.Header.Set("X-Riot-Token", s.apiKey)

	res, err := s.client.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		err = errors.New("Unsuccessful status code: " + fmt.Sprint(res.StatusCode))
		return
	}

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return
	}

	return
}

func (s LolApi) GetAccountByRiotId(gameName string, tagLine string) (acc account.Account, err error) {
	logger.Debug(s, "Getting account by riot id")

	url := AMERICAS_URL +
		"/riot/account/v1/accounts/by-riot-id/" +
		gameName + "/" + tagLine

	err = s.fetch("GET", url, &acc)
	if err != nil {
		logger.Error(s, "Cannot get account by riot id", err)
	}

	return
}

func (s LolApi) GetLeaguesBySummonerId(id string) (lgs []lol.LeagueItem, err error) {
	logger.Debug(s, "Getting leagues by summoner id")

	url := BR1_LOL_URL +
		"/league/v4/entries/by-summoner/" +
		id

	err = s.fetch("GET", url, &lgs)
	if err != nil {
		logger.Error(s, "Cannot get leagues by summoner id", err)
	}

	return
}

func (s LolApi) GetSummonerByPuuid(puuid string) (smm lol.Summoner, err error) {
	logger.Debug(s, "Getting summoner by puuid")

	url := BR1_LOL_URL +
		"/summoner/v4/summoners/by-puuid/" +
		puuid

	err = s.fetch("GET", url, &smm)
	if err != nil {
		logger.Error(s, "Cannot get summoner by puuid", err)
	}

	return
}

func (s LolApi) GetMatchIdListByPuuid(puuid string) (mil []string, err error) {
	logger.Debug(s, "Getting match id list by puuid")

	url := AMERICAS_URL +
		"/lol/match/v5/matches/by-puuid/" +
		puuid +
		"/ids"

	err = s.fetch("GET", url, &mil)
	if err != nil {
		logger.Error(s, "Cannot get match id list by puuid", err)
	}

	return
}

func (s LolApi) GetMatchByMatchId(matchId string) (me lol.Match, err error) {
	logger.Debug(s, "Getting match by match id")

	url := AMERICAS_URL +
		"/lol/match/v5/matches/" +
		matchId

	err = s.fetch("GET", url, &me)
	if err != nil {
		logger.Error(s, "Cannot get match by match id", err)
	}

	return
}
