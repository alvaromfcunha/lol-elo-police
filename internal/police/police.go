package police

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/alvaromfcunha/lol-elo-police/internal/model"
	"github.com/alvaromfcunha/lol-elo-police/pkg/lol"
	"github.com/alvaromfcunha/lol-elo-police/pkg/wpp"
	"gorm.io/gorm"
)

type Police struct {
	Db        *gorm.DB
	LolApi    lol.LolApi
	WppClient wpp.WhatsappClient
	Interval  time.Duration
	GroupUser string
}

type PlayerUpdate struct {
	LeagueEntry lol.LeagueEntry
	Player      model.Player
}

func (p Police) Start() {
	messagesFile := "messages.txt"
	tmpl, err := template.New(messagesFile).ParseFiles(messagesFile)
	if err != nil {
		panic("can't load " + messagesFile + " template file")
	}

	for {
		time.Sleep(p.Interval)
		p.PatrolJob(tmpl)
	}
}

func (p Police) PatrolJob(tmpl *template.Template) {
	var players []model.Player
	p.Db.Find(&players)

	for _, player := range players {
		go func(player model.Player) {
			leagues, err := p.LolApi.GetLeaguesBySummonerId(player.SummonerId)
			if err != nil {
				return
			}

			var solo *lol.LeagueEntry
			for _, league := range leagues {
				if league.QueueType == lol.QueueType(lol.Solo) {
					solo = &league
					break
				}
			}

			if solo.LeaguePoints != player.LeaguePoints ||
				solo.Wins != player.Wins ||
				solo.Losses != player.Losses ||
				solo.Rank != player.Rank ||
				solo.Tier != player.Tier {

				update := PlayerUpdate{
					LeagueEntry: *solo,
					Player:      player,
				}

				var textBuf bytes.Buffer
				err = tmpl.Execute(&textBuf, update)
				if err != nil {
					fmt.Println(err)
				}

				p.WppClient.SendMessageToGroup(textBuf.String(), p.GroupUser)

				player.LeaguePoints = solo.LeaguePoints
				player.Wins = solo.Wins
				player.Losses = solo.Losses
				player.Rank = solo.Rank
				player.Tier = solo.Tier

				p.Db.Save(player)
			}
		}(player)
	}
}
