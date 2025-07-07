package utils

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("1"),
		tgbotapi.NewKeyboardButton("2"),
		tgbotapi.NewKeyboardButton("3"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("4"),
		tgbotapi.NewKeyboardButton("5"),
		tgbotapi.NewKeyboardButton("6"),
	),
)

var StartKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Мероприятия"),
		tgbotapi.NewKeyboardButton("Вступить"),
		tgbotapi.NewKeyboardButton("Зарегистрироваться"),
	),
)

var PreApplicationKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Заполнить анкету"),
		tgbotapi.NewKeyboardButton("Вернуться"),
	),
)

func GenerateApplicationKeyboard(userID int64, selected []string) tgbotapi.InlineKeyboardMarkup {
	options := []string{"Организаторы", "Медиа", "Редакторы", "Дизайнеры"}
	selectedMap := make(map[string]bool)
	for _, s := range selected {
		selectedMap[s] = true
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, opt := range options {
		prefix := "⬜️"
		if selectedMap[opt] {
			prefix = "✅"
		}
		btn := tgbotapi.NewInlineKeyboardButtonData(prefix+" "+opt, "toggle_"+opt)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}
	done := tgbotapi.NewInlineKeyboardButtonData("✅ Готово", "done")
	cancel := tgbotapi.NewInlineKeyboardButtonData("ОТМЕНИТЬ", "cancel")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(done, cancel))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
