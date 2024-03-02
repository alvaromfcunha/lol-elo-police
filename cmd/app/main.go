package main

import (
	"os"
	"strconv"
	"time"

	playerHandler "github.com/alvaromfcunha/lol-elo-police/internal/handler/player"
	"github.com/alvaromfcunha/lol-elo-police/internal/model"
	"github.com/alvaromfcunha/lol-elo-police/internal/police"
	"github.com/alvaromfcunha/lol-elo-police/pkg/lol"
	"github.com/alvaromfcunha/lol-elo-police/pkg/wpp"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	riotApiKey := os.Getenv("RIOT_API_KEY")
	wppGroupUser := os.Getenv("WPP_GROUP_USER")
	policeIntervalMinutes := os.Getenv("POLICE_INTERVAL_MINUTES")

	minutes, err := strconv.Atoi(policeIntervalMinutes)
	if err != nil {
		panic("invalid POLICE_INTERVAL_MINUTES dotenv var")
	}

	db, err := gorm.Open(sqlite.Open("db/app.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.Player{})

	app := fiber.New()

	wpp := wpp.WhatsappClient{}
	err = wpp.Init()
	if err != nil {
		panic(err)
	}

	lol := lol.LolApi{
		ApiKey: riotApiKey,
	}
	pol := police.Police{
		Interval:  time.Duration(minutes) * time.Minute,
		Db:        db,
		LolApi:    lol,
		WppClient: wpp,
		GroupUser: wppGroupUser,
	}
	ph := playerHandler.PlayerHandler{
		Db:     db,
		LolApi: lol,
	}

	app.Post("/player", ph.CreatePlayer)
	app.Get("/player", ph.ReadAllPlayers)

	go pol.Start()
	app.Listen(":3000")
}
