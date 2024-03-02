package police

import (
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

func (p Police) Start() {
	for {
		time.Sleep(p.Interval)
		p.PatrolJob()
	}
}

func (p Police) PatrolJob() {
	var players []model.Player
	p.Db.Find(&players)

	for _, player := range players {
		leagues, err := p.LolApi.GetLeaguesBySummonerId(player.SummonerId)
		if err != nil {
			continue
		}

		var solo *lol.LeagueEntry
		for _, league := range leagues {
			if league.QueueType == lol.QueueType(lol.Solo) {
				solo = &league
			}
		}

		if solo.LeaguePoints != player.LeaguePoints ||
			solo.Wins != player.Wins ||
			solo.Losses != player.Losses ||
			solo.Rank != player.Rank ||
			solo.Tier != player.Tier {

			text := player.GameName + "#" + player.TagLine + " mudou!"
			p.WppClient.SendMessageToGroup(text, p.GroupUser)

			player.LeaguePoints = solo.LeaguePoints
			player.Wins = solo.Wins
			player.Losses = solo.Losses
			player.Rank = solo.Rank
			player.Tier = solo.Tier

			p.Db.Save(player)
		}
	}
}
