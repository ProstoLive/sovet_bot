package models

import (
	"time"

	"github.com/lib/pq"
)

type User struct {
	ID         int64  `db:"id"`
	TelegramID int64  `db:"telegram_id"`
	FirstName  string `db:"first_name"`
	LastName   string `db:"last_name"`
	Username   string `db:"username"`
	CreatedAt  string `db:"created_at"`
}

type ApplicationState struct {
	UserID     int64          `db:"user_id"`
	Directions pq.StringArray `db:"directions"`
	UpdatedAt  time.Time      `db:"updated_at"`
}
