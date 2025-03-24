package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID           uuid.UUID `json:"id"`
	CalendarID   uuid.UUID `json:"calendar_id"`
	Description  string    `json:"description"`
	Location     string    `json:"location"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	NotifyBefore time.Duration `json:"notify_before"`
	CreatedAt    time.Time `json:"created_at"`
}
