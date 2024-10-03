package service

import "go.mau.fi/whatsmeow"

type IWhatsappService interface {
	SendMessageToGroup(text string, group string) (resp whatsmeow.SendResponse, err error)
}
