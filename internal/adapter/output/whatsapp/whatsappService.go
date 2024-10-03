package whatsapp

import (
	"context"
	"errors"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

type WhatsappService struct {
	Client *whatsmeow.Client
}

func NewWhatsappService(client *whatsmeow.Client) WhatsappService {
	return WhatsappService{
		Client: client,
	}
}

func (s WhatsappService) SendMessageToGroup(text string, group string) (resp whatsmeow.SendResponse, err error) {
	if !s.Client.IsConnected() {
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

	resp, err = s.Client.SendMessage(
		context.Background(),
		jid,
		msg,
	)

	return
}
