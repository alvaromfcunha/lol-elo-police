package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alvaromfcunha/lol-elo-police/pkg/wpp"
	"go.mau.fi/whatsmeow/types/events"
)

func main() {
	c, err := wpp.GetClient()
	if err != nil {
		return
	}

	if c.Store.ID != nil {
		err = errors.New("device registered in store already")
		panic(err)
	}

	qrChan, _ := c.GetQRChannel(context.Background())
	err = c.Connect()
	if err != nil {
		panic(err)
	}
	c.AddEventHandler(func(evt interface{}) {
		switch evt.(type) {
		case *events.AppStateSyncComplete:
			fmt.Println("done!")
			c.Disconnect()
		}
	})

	for evt := range qrChan {
		if evt.Event == "code" {
			fmt.Println(evt.Code)
		}
	}

	e := make(chan os.Signal)
	signal.Notify(e, os.Interrupt, syscall.SIGTERM)
	<-e

	c.Disconnect()
}
