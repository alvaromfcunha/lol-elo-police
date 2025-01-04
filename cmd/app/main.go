package main

import (
	"text/template"

	"github.com/alvaromfcunha/lol-elo-police/cmd/app/controller"
	"github.com/alvaromfcunha/lol-elo-police/cmd/app/util"
	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"

	"database/sql"

	"github.com/joho/godotenv"
)

func main() {
	// Env.
	err := godotenv.Load("./infrastructure/config/.env")
	if err != nil {
		panic("cannot load env")
	}

	// Database.
	db, err := sql.Open("sqlite3", "file:./infrastructure/database/app.db")
	if err != nil {
		panic("cannot connect to database")
	}

	// Template engine.
	messageTemplates, err := template.New("template").Funcs(util.TemplateFuncMap).ParseFiles("./infrastructure/template/messages.txt")
	if err != nil {
		panic("cannot load template file: " + err.Error())
	}

	// Whatsapp/Whatsmeow.
	container, err := sqlstore.New("sqlite3", "file:./infrastructure/database/whatsapp.db?_foreign_keys=on", nil)
	if err != nil {
		panic("cannot load whatsmeow store from sqlite file: " + err.Error())
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic("cannot retrieve registered device on whatsmeow store")
	}
	whatsappClient := whatsmeow.NewClient(deviceStore, nil)
	if whatsappClient.Store.ID == nil {
		panic("cannot load registered device on whatsmeow store")
	}
	err = whatsappClient.Connect()
	if err != nil {
		panic("cannot connect whatsmeow client")
	}

	// Scheduler.
	cron := cron.New()

	cron.AddFunc("*/1 * * * *", controller.PatrolRoutineController(db, messageTemplates, whatsappClient))

	cron.Start()

	// API.
	api := fiber.New()

	api.Post("/players", controller.CreatePlayerController(db))
	api.Get("/players", controller.GetPlayersController(db))
	api.Post("/execute/patrol", controller.ExecutePatrolController(db, messageTemplates, whatsappClient))

	api.Listen(":3000") // blocking
}
