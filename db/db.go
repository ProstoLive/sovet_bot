package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sovet_bot/db/models"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func Connect() {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
}

func CheckIfNotExists(user *tgbotapi.User) {
	var dbUser models.User
	err := DB.Get(&dbUser, "SELECT * FROM users WHERE telegram_id=$1", user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = DB.Exec(
				`INSERT INTO users (telegram_id, username, first_name, last_name, created_at) VALUES ($1, $2, $3, $4, $5)`,
				user.ID, user.UserName, user.FirstName, user.LastName, time.Now(),
			)
			if err != nil {
				log.Printf("New user insertion error: %v", err)
			} else {
				log.Printf("New user %v has been added to database!", user.UserName)
			}
		} else {
			log.Fatal(err)
		}
	} else {
		log.Printf("User %v already exists in database", user.UserName)
	}

}

// Получение текущего состояния или создание нового
func GetOrCreateState(userID int64) (*models.ApplicationState, error) {
	var state models.ApplicationState
	err := DB.Get(&state, "SELECT * FROM application_states WHERE user_id = $1", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err := DB.Exec("INSERT INTO application_states (user_id, directions) VALUES ($1, $2)", userID, pq.StringArray([]string{}))
			if err != nil {
				return nil, err
			}
			state = models.ApplicationState{UserID: userID, Directions: []string{}, UpdatedAt: time.Now()}
			return &state, nil
		}
		return nil, err
	}

	if state.Directions == nil {
		state.Directions = []string{}
	}

	return &state, nil
}

// Добавление или удаление направления (toggle)
func ToggleDirection(userID int64, direction string) error {
	state, err := GetOrCreateState(userID)
	if err != nil {
		return err
	}

	updated := false
	newDirs := []string{}
	found := false

	for _, d := range state.Directions {
		if d == direction {
			found = true // удаляем
			updated = true
			continue
		}
		newDirs = append(newDirs, d)
	}
	if !found {
		newDirs = append(newDirs, direction)
		updated = true
	}

	if updated {
		_, err = DB.Exec("UPDATE application_states SET directions = $1, updated_at = $2 WHERE user_id = $3", pq.StringArray(newDirs), time.Now(), userID)
	}

	return err
}

// Получить финальный список направлений
func GetSelectedDirections(userID int64) ([]string, error) {
	var directions []string
	err := DB.Get(&directions, "SELECT directions FROM application_states WHERE user_id = $1", userID)
	return directions, err
}

// Очистить после завершения
func ClearState(userID int64) error {
	_, err := DB.Exec("DELETE FROM application_states WHERE user_id = $1", userID)
	return err
}

func CheckIfFillingApplication(userID int64) bool {
	var directions []string
	err := DB.Get(&directions, "SELECT directions FROM application_states WHERE user_id = $1", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}
	}
	return true
}
