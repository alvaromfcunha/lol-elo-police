package main

import (
	"time"

	playerHandler "github.com/alvaromfcunha/lol-elo-police/internal/handler/player"
	"github.com/alvaromfcunha/lol-elo-police/internal/model"
	"github.com/alvaromfcunha/lol-elo-police/pkg/lol"
	"github.com/alvaromfcunha/lol-elo-police/pkg/police"
	"github.com/alvaromfcunha/lol-elo-police/pkg/wpp"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("db/app.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.Player{})

	wpp := wpp.WhatsappClient{}
	err = wpp.Init()
	if err != nil {
		panic(err)
	}

	lol := lol.LolApi{
		ApiKey: "RGAPI-b11d5bf3-b56d-4856-b50f-d685f7495360",
	}

	ph := playerHandler.PlayerHandler{
		Db:     db,
		LolApi: lol,
	}

	app := fiber.New()
	app.Post("/player", ph.CreatePlayer)
	app.Get("/player", ph.ReadAllPlayers)
	go app.Listen(":3000")

	pol := police.Police{
		Interval:  5 * time.Minute,
		Db:        db,
		LolApi:    lol,
		WppClient: wpp,
		GroupUser: "553599945538-1596561080",
	}
	pol.Start()

	<-make(chan bool)
}
