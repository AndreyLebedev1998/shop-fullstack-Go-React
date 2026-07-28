package interfaces

import (
	"net/smtp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TgBot struct {
	ChatId int64
	Bot    *tgbotapi.BotAPI
}

type EmailData struct {
	Auth smtp.Auth
	From string
	To   []string
}

type Notifier interface {
	Send(msg []byte) error
}

func (t TgBot) Send(message []byte) error {
	str := string(message)
	msg := tgbotapi.NewMessage(t.ChatId, str)
	_, err := t.Bot.Send(msg)
	return err
}

func (e EmailData) Send(message []byte) error {
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		e.Auth,
		e.From,
		e.To,
		message,
	)
	return err
}
