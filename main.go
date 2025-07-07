package main

import (
	"log"
	"os"
	"sovet_bot/db"
	"sovet_bot/utils"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	db.Connect()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.CallbackQuery != nil {
			data := update.CallbackQuery.Data
			userID := update.CallbackQuery.From.ID
			chatID := update.CallbackQuery.Message.Chat.ID

			if strings.HasPrefix(data, "toggle_") {
				dir := strings.TrimPrefix(data, "toggle_")
				err := db.ToggleDirection(userID, dir)
				if err != nil {
					log.Println("Toggle error:", err)
				}

				state, _ := db.GetOrCreateState(userID)

				edit := tgbotapi.NewEditMessageTextAndMarkup(
					chatID,
					update.CallbackQuery.Message.MessageID,
					"Выберите направления:",
					utils.GenerateApplicationKeyboard(userID, state.Directions),
				)
				bot.Send(edit)
				bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
			} else if data == "done" {
				//Удаление inline-клавиатуры
				edit := tgbotapi.NewEditMessageReplyMarkup(
					chatID,
					update.CallbackQuery.Message.MessageID,
					tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{}},
				)
				bot.Send(edit)

				dirs, _ := db.GetSelectedDirections(userID)
				bot.Send(tgbotapi.NewMessage(chatID, "Вы выбрали: "+strings.Join(dirs, ", ")))
				db.ClearState(userID)
				bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))

			} else if data == "cancel" {
				//Удаление inline-клавиатуры
				edit := tgbotapi.NewEditMessageReplyMarkup(
					chatID,
					update.CallbackQuery.Message.MessageID,
					tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{}},
				)
				bot.Send(edit)

				db.ClearState(userID)
				msg := tgbotapi.NewMessage(chatID, utils.StartMessage)
				msg.ParseMode = "HTML"
				msg.ReplyMarkup = utils.StartKeyboard
				bot.Send(msg)
				bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
			}
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if db.CheckIfFillingApplication(update.Message.From.ID) {
			msg.Text = utils.ApplicationFillingMessage
			msg.ParseMode = "HTML"
			bot.Send(msg)
			continue
		}

		if update.Message == nil {
			continue
		}

		msg.Text = update.Message.Text

		switch update.Message.Text {
		case "/start":
			db.CheckIfNotExists(update.Message.From)
			msg.ParseMode = "HTML"
			msg.Text = utils.StartMessage
			msg.ReplyMarkup = utils.StartKeyboard
		case "open":
			msg.ReplyMarkup = utils.StartKeyboard
		case "close":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		case "default":
			msg.ParseMode = "HTML"
			msg.Text = utils.StartMessage
		case "Вступить":
			msg.Text = utils.SovetMessage
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			msg.ReplyMarkup = utils.PreApplicationKeyboard

			//msg.ReplyMarkup = utils.GenerateApplicationKeyboard(update.Message.From.ID, state.Directions)
		case "Заполнить анкету":
			db.GetOrCreateState(update.Message.From.ID)

		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
