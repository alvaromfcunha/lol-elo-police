package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
)

func main() {
	container, err := sqlstore.New("sqlite3", "file:infrastructure/database/whatsapp.db?_foreign_keys=on", nil)
	if err != nil {
		panic("cannot load whatsmeow store from sqlite file")
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic("cannot retrieve registered device on whatsmeow store")
	}
	client := whatsmeow.NewClient(deviceStore, nil)
	if client.Store.ID != nil {
		panic("device already registered")
	}

	qrChan, _ := client.GetQRChannel(context.Background())
	err = client.Connect()
	if err != nil {
		panic(err)
	}
	client.AddEventHandler(func(evt interface{}) {
		switch evt.(type) {
		case *events.AppStateSyncComplete:
			fmt.Println("done!")
			client.Disconnect()
		}
	})

	for evt := range qrChan {
		if evt.Event == "code" {
			fmt.Println(evt.Code)
		}
	}

	e := make(chan os.Signal, 1)
	signal.Notify(e, os.Interrupt, syscall.SIGTERM)
	<-e

	client.Disconnect()
}
