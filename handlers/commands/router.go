package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler struct {
	Bot *tgbotapi.BotAPI
	Message *tgbotapi.Message
	//Остальные Repo
}

type CallbackHandler struct {
	Bot *tgbotapi.BotAPI
	Callback *tgbotapi.CallbackConfig
	//Остальные Repo
}


func (h *CommandHandler) Execute() error {
	switch h.Message.Text {
	case "/start":
		return SendStartMessage(h.Bot, h.Message)
	default:
		return nil
	}
}
