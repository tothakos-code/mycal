package models

import (
	"time"

	"github.com/google/uuid"
)

type Invitation struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	EventID  uuid.UUID `json:"event_id"`
	Status   string    `json:"status"` // "pending", "accepted", "declined"
	CreatedAt time.Time `json:"created_at"`
}
