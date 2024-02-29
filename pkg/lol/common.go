package lol

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	AMERICAS_RIOT_URL = "https://americas.api.riotgames.com/riot"
	BR1_LOL_URL       = "https://br1.api.riotgames.com/lol"
)

type LolApi struct {
	ApiKey     string
	HttpClient http.Client
}

func DoRequest[T any](l LolApi, method string, url string, t *T) (err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Set("X-Riot-Token", l.ApiKey)

	res, err := l.HttpClient.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		err = errors.New("unsuccessful status code: " + fmt.Sprint(res.StatusCode))
		return
	}

	err = json.NewDecoder(res.Body).Decode(&t)
	if err != nil {
		return
	}

	return
}
