package callback

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackHandler struct {
	Bot      *tgbotapi.BotAPI
	Callback *tgbotapi.CallbackQuery
}

func (h *CallbackHandler) Execute() error {
	if h.Callback.Data == "" {
		return errors.New("expected CallbackQuery")
	}
	data := h.Callback.Data
	// Логика обработки callback
	msg := tgbotapi.NewMessage(h.Callback.Message.Chat.ID, "Callback: "+data)
	_, err := h.Bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
