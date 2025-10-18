package services

import (
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	jsonrpc "github.com/gorilla/rpc/json"
)

type SendMessageArgs struct {
	Aboba string `json:"aboba"`
}

type YaFormsService struct {
	Bot *tgbotapi.BotAPI
}

func (y *YaFormsService) handleYaForms(r *http.Request, args *SendMessageArgs) error {
	msg := tgbotapi.NewMessage(args.ChatID, args.Aboba)
	_, err := y.Bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func CreateNewRpc(bot *tgbotapi.BotAPI) {
	s := rpc.NewServer()
	s.RegisterCodec(jsonrpc.NewCodec(), "application/json")
	yfs := &YaFormsService{Bot: bot}
	s.RegisterService(yfs, "Telegram")

	router := mux.NewRouter()

	router.Handle("/rpc", s)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
