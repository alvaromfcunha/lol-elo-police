package whatsapp

import (
	"context"
	"errors"

	"github.com/alvaromfcunha/lol-elo-police/internal/adapter/output/logger"
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
	logger.Debug(s, "Sending message to group")

	if !s.Client.IsConnected() {
		err = errors.New("not connected")
		logger.Error(s, "Cannot send message to group", err)
		return
	}

	msg := &proto.Message{
		Conversation: &text,
	}

	jid := types.JID{
		User:   group,
		Server: types.GroupServer,
	}

	s.Client.SendPresence(types.PresenceAvailable)
	s.Client.SendChatPresence(jid, types.ChatPresenceComposing, types.ChatPresenceMediaText)

	resp, err = s.Client.SendMessage(
		context.Background(),
		jid,
		msg,
	)
	if err != nil {
		logger.Error(s, "Cannot send message to group", err)
	}

	s.Client.SendChatPresence(jid, types.ChatPresencePaused, types.ChatPresenceMediaText)
	s.Client.SendPresence(types.PresenceUnavailable)

	return
}
