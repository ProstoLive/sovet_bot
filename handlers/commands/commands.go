package commands

import (
	"log"
	"sovet_bot/db"
	"sovet_bot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendStartMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) error {

	db.CheckIfNotExists(message.From)
	msg := tgbotapi.NewMessage(0, "")

	msg.ChatID = message.Chat.ID
	msg.ReplyMarkup = utils.StartKeyboard
	msg.Text = utils.StartMessage

	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
	return nil
}
