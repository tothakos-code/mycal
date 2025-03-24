package handler

import (
	"github.com/google/uuid"
	"golang-postgresql-auth-template/internal/app/repository"
	"golang-postgresql-auth-template/internal/models"
	"log"
	"net/http"
	"time"
	"encoding/json"
)

type EventHandler struct {
	eventRepo *repository.EventRepo
}

func NewEventHandler(eventRepo *repository.EventRepo) *EventHandler {
	return &EventHandler{eventRepo: eventRepo}
}

func (e *EventHandler) HandleCreateEvent() http.HandlerFunc {
	type request struct {
		CalendarID   uuid.UUID `json:"calendar_id" validate:"required,uuid"`
		Description  string    `json:"description"`
		Location     string    `json:"location"`
		StartDate    time.Time `json:"start_date"`
		EndDate      time.Time `json:"end_date"`
		StartTime    time.Time `json:"start_time"`
		EndTime      time.Time `json:"end_time"`
		NotifyBefore time.Duration `json:"notify_before"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
    form, err := decodeJsonReq[request](r)
		if err != nil {
			log.Println("# Error decoding request", err)
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		event := models.Event{
			CalendarID:   form.CalendarID,
			Description:  form.Description,
			Location:     form.Location,
			StartDate:    form.StartDate,
			EndDate:      form.EndDate,
			StartTime:    form.StartTime,
			EndTime:      form.EndTime,
			NotifyBefore: time.Duration(form.NotifyBefore) * time.Minute,
			CreatedAt:    time.Now(),
		}

		err = e.eventRepo.CreateEvent(ctx, event)
		if err != nil {
			log.Println("# Error creating event:", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (e *EventHandler) HandleListPublicEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		events, err := e.eventRepo.ListPublicEvents(ctx)
		if err != nil {
			log.Println("# Error listing public events:", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	}
}

func (e *EventHandler) HandleListEventsByCalendars() http.HandlerFunc {
	type request struct {
		CalendarIDs []string `json:"calendar_ids" validate:"required,dive,uuid"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var form request
		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		events, err := e.eventRepo.ListEventsByCalendars(ctx, form.CalendarIDs)
		if err != nil {
			log.Println("# Error listing events:", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	}
}
