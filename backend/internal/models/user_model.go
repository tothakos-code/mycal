package models

import (
	"time"

	"github.com/google/uuid"
)


type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Username     string    `json:"username"`
	FirstName    string    `json:"firstname"`
	SurName      string    `json:"surname"`
	CreatedAt    time.Time `json:"created_at"`
}
