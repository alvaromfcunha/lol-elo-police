package wpp

import (
	"context"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
)

type WhatsappClient struct {
	Client     *whatsmeow.Client
	LastQrCode string
}

func GetClient() (client *whatsmeow.Client, err error) {
	container, err := sqlstore.New("sqlite3", "file:db/wpp.db?_foreign_keys=on", nil)
	if err != nil {
		return
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return
	}

	client = whatsmeow.NewClient(deviceStore, nil)
	return
}

func (w *WhatsappClient) Init() (err error) {
	w.Client, err = GetClient()
	if err != nil {
		return
	}

	if w.Client.Store.ID == nil {
		err = errors.New("no device registered in store")
		return
	}

	err = w.Client.Connect()
	if err != nil {
		return
	}

	return
}

func (w WhatsappClient) SendMessageToGroup(text string, group string) (resp whatsmeow.SendResponse, err error) {
	if !w.Client.IsConnected() {
		err = errors.New("not connected")
		return
	}

	msg := &proto.Message{
		Conversation: &text,
	}

	jid := types.JID{
		User:   group,
		Server: types.GroupServer,
	}

	resp, err = w.Client.SendMessage(
		context.Background(),
		jid,
		msg,
	)

	return
}
