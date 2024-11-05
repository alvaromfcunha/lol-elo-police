package main

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/entity"
	"github.com/alvaromfcunha/lol-elo-police/internal/domain/service"
)

func main() {
	messageTemplates, err := template.New("./template.txt").Funcs(template.FuncMap{"sub": func(i int, sub int) int { return i - sub }}).ParseFiles("./template.txt")
	if err != nil {
		panic(err)
	}

	match := entity.NewMatch(
		"123",
		420,
		0,
		0,
		0,
	)

	mp := entity.NewMatchParticipant(
		match,
		entity.NewPlayer(
			"123",
			"123",
			"Player1",
			"BR1",
		),
		"Aatrox",
		"Support",
		9,
		0,
		1,
		true,
	)

	li := &lol.LeagueItem{
		QueueType:    "string",
		SummonerName: "string",
		HotStreak:    true,
		Wins:         1,
		Veteran:      true,
		Losses:       1,
		FreshBlood:   true,
		Inactive:     true,
		Tier:         "string",
		Rank:         "string",
		SummonerID:   "string",
		LeaguePoints: 1,
	}

	data := service.NewRankedMatchData{
		Match: match,
		MatchParticipantsWithLeagueItem: []service.MatchParticipantWithLeagueItem{
			{
				MatchParticipant: mp,
				LeagueItem:       li,
			},
			{
				MatchParticipant: mp,
				LeagueItem:       li,
			},
			{
				MatchParticipant: mp,
				LeagueItem:       li,
			},
			{
				MatchParticipant: mp,
				LeagueItem:       li,
			},
			{
				MatchParticipant: mp,
				LeagueItem:       li,
			},
		},
	}

	var buf bytes.Buffer
	err = messageTemplates.ExecuteTemplate(&buf, "NewRankedMatch", data)
	if err != nil {
		panic(err)
	}

	fmt.Println(buf.String())
}
