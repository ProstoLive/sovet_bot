package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gorilla/mux"
)

type FormArgs struct {
	Fio       string `json:"fio"`
	Group     string `json:"group"`
	Link      string `json:"link"`
	Spheres   string `json:"spheres"`
	DifSphere string `json:"dif_sphere,omitempty"`
	Hobbies   string `json:"hobbies"`
}

type SendMessageReply struct {
	Result string `json:"result"`
}

type YaFormsService struct {
	Bot *tgbotapi.BotAPI
}

func (y *YaFormsService) HandleYaForms(w http.ResponseWriter, r *http.Request) {
	var args FormArgs

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &args); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	adminChatId, err := strconv.ParseInt(os.Getenv("ADMIN_CHAT_ID"), 10, 64)
	if err != nil {
		http.Error(w, "Internal server", http.StatusInternalServerError)
	}

	difSphere := ""
	if args.DifSphere != "" {
		difSphere = "\nДругое: " + args.DifSphere
	}

	sendMsg := fmt.Sprintf(
		`Новая заявка в студсовет! 
		ФИО: %s\n\nГруппа: %s\n\nСсылка: %s\n\nНаправления: %s %s\n\nО себе: %s`,
		args.Fio, args.Group, args.Link, args.Spheres, difSphere, args.Hobbies,
	)

	msg := tgbotapi.NewMessage(adminChatId, sendMsg)
	_, err = y.Bot.Send(msg)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	reply := SendMessageReply{Result: "Message sent"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reply)
}

func CreateNewYaFormsServer(bot *tgbotapi.BotAPI) {
	yfs := &YaFormsService{Bot: bot}

	router := mux.NewRouter()

	router.HandleFunc("/yaforms", yfs.HandleYaForms).Methods("POST")

	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
