package main

import (
	"log"
	"os"
	"sovet_bot/db"
	"sovet_bot/handlers"
	"sovet_bot/services"

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

	go services.CreateNewRpc(bot)

	updates := bot.GetUpdatesChan(u)

	factory := &handlers.HandlerFactory{Bot: bot}
	for update := range updates {
		handler := factory.GetHandler(&update)
		if handler != nil {
			if err := handler.Execute(); err != nil {
				// обработка ошибок, логгирование и пр.
			}
		} else {
			// нет обработчика, можно залогировать или проигнорировать
		}
	}
}
