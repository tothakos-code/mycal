package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	EventID   uuid.UUID `json:"event_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	ShownAt   time.Time `json:"shown_at"`
	CreatedAt time.Time `json:"created_at"`
}
