package handlers

import (
	"sovet_bot/handlers/callback"
	"sovet_bot/handlers/commands"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler interface {
	Execute() error
}

type HandlerFactory struct {
	Bot *tgbotapi.BotAPI
	// тут можно добавить репозитории, зависимости
}

func (f *HandlerFactory) GetHandler(update *tgbotapi.Update) Handler {
	switch {
	case update.Message != nil:
		return &commands.CommandHandler{Bot: f.Bot, Message: update.Message}
	case update.CallbackQuery != nil:
		return &callback.CallbackHandler{Bot: f.Bot, Callback: update.CallbackQuery}
	}
	return nil
}
