package handler

import (
	"golang-postgresql-auth-template/internal/app/repository"
	"golang-postgresql-auth-template/internal/models"
  "golang-postgresql-auth-template/internal/app/service"
	"log"
	"net/http"
	"time"
	"encoding/json"
  "github.com/google/uuid"
)

type EventHandler struct {
	eventRepo *repository.EventRepo
	UserRepo *repository.UserRepo
}

func NewEventHandler(eventRepo *repository.EventRepo, UserRepo *repository.UserRepo) *EventHandler {
	return &EventHandler{eventRepo: eventRepo, UserRepo: UserRepo}
}

func (e *EventHandler) HandleCreateEvent() http.HandlerFunc {
	type request struct {
		UserID       uuid.UUID `json:"user_id"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		Location     string    `json:"location"`
		Start        time.Time `json:"start"`
		Finish       time.Time `json:"finish"`
		Public       bool      `json:"is_public"`
		NotifyBefore uint16    `json:"notify_before"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
    form, err := decodeJsonReq[request](r)
		if err != nil {
			log.Println("# Error decoding request", err)
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}
    user_id, ok := service.GetUserIDFromContext(r)
    if !ok {
      log.Println("Error: no user in context", err)
      http.Error(w, "Invalid Request", http.StatusBadRequest)
      return
    }

		event := models.Event{
      UserID:       user_id,
			Title:        form.Title,
			Description:  form.Description,
			Location:     form.Location,
			Start:        form.Start,
			Finish:       form.Finish,
			Public:       form.Public,
			NotifyBefore: form.NotifyBefore,
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
		publicEvents, err := e.eventRepo.ListPublicEvents(ctx)
		if err != nil {
			log.Println("Error listing public events:", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
    userId, ok := service.GetUserIDFromContext(r)
    if !ok {
      log.Println("Error user_id not fount in context:", err)
      http.Error(w, "Server Error", http.StatusInternalServerError)
      return
    }
		privateEvents, err := e.eventRepo.ListEventsByUserID(ctx, userId, false)
		if err != nil {
			log.Println("Error listing user's private events:", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
    var allEvents []models.Event
    allEvents = append(publicEvents, privateEvents...)
		w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(allEvents)

	}
}

// func (e *EventHandler) HandleListEventsByUser() http.HandlerFunc {
// 	type request struct {
// 		CalendarIDs []string `json:"calendar_ids" validate:"required,dive,uuid"`
// 	}
//
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()
// 		var form request
// 		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
// 			http.Error(w, "Invalid Request", http.StatusBadRequest)
// 			return
// 		}
//
// 		events, err := e.eventRepo.ListEventsByUserID(ctx, form.CalendarIDs)
// 		if err != nil {
// 			log.Println("# Error listing events:", err)
// 			http.Error(w, "Server Error", http.StatusInternalServerError)
// 			return
// 		}
//
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(events)
// 	}
// }
