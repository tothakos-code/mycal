package handler

import (
	"encoding/json"
	"golang-postgresql-auth-template/internal/app/repository"
	"log"
	"net/http"
  "github.com/google/uuid"
)

type CalendarHandler struct {
	calendarRepo *repository.CalendarRepo
}

func NewCalendarHandler(calendarRepo *repository.CalendarRepo) *CalendarHandler {
	return &CalendarHandler{calendarRepo: calendarRepo}
}

func (c *CalendarHandler) HandleCreateCalendar() http.HandlerFunc {
	type request struct {
		UserID   uuid.UUID `json:"user_id" validate:"required,uuid"`
		IsPublic bool   `json:"is_public"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var form request
		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		err := c.calendarRepo.CreateCalendar(ctx, form.UserID, form.IsPublic)
		if err != nil {
			log.Println("# Error creating calendar:", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
