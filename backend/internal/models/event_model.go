package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Location     string    `json:"location"`
	Start        time.Time `json:"start"`
	Finish       time.Time `json:"finish"`
	NotifyBefore uint16    `json:"notify_before"`
  Public       bool      `json:"is_public"`
	CreatedAt    time.Time `json:"created_at"`
}
