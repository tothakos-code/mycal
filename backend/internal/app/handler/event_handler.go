package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"golang-postgresql-auth-template/internal/app/repository"
	"golang-postgresql-auth-template/internal/app/service"
	"golang-postgresql-auth-template/internal/models"
	"log"
	"net/http"
	"strings"
	"time"
)

type EventHandler struct {
	eventRepo *repository.EventRepo
	UserRepo  *repository.UserRepo
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

// HandleEditEvent handles updating an existing event.
func (e *EventHandler) HandleUpdateEvent() http.HandlerFunc {
	type request struct {
		Title        *string    `json:"title"`         // Use pointers to allow for optional fields
		Description  *string    `json:"description"`   // Use pointers to allow for optional fields
		Location     *string    `json:"location"`      // Use pointers to allow for optional fields
		Start        *time.Time `json:"start"`         // Use pointers to allow for optional fields
		Finish       *time.Time `json:"finish"`        // Use pointers to allow for optional fields
		Public       *bool      `json:"is_public"`     // Use pointers to allow for optional fields
		NotifyBefore *uint16    `json:"notify_before"` // Use pointers to allow for optional fields
	}

	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		log.Println(parts)
		if len(parts) != 3 || parts[2] == "" {
			log.Println("Invalid event ID in path")
			http.Error(w, "Invalid event ID in path", http.StatusBadRequest)
			return
		}

		idStr := parts[2] // Get the "1234" part
		eventID, err := uuid.Parse(idStr)
		if err != nil {
			log.Println("Invalid UUID")
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		// Decode the request body into the form struct
		ctx := r.Context()
		form, err := decodeJsonReq[request](r)
		if err != nil {
			log.Println("# Error decoding request", err)
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		// Retrieve the user ID from context
		user_id, ok := service.GetUserIDFromContext(r)
		if !ok {
			log.Println("Error: no user in context", err)
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		// Fetch the existing event from the repository
		event, err := e.eventRepo.GetEventByID(ctx, eventID)
		if err != nil {
			log.Println("# Error fetching event:", err)
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}

		// Ensure the event belongs to the user (if applicable)
		if event.UserID != user_id {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Update fields in the event if they are provided in the request
		if form.Title != nil {
			event.Title = *form.Title
		}
		if form.Description != nil {
			event.Description = *form.Description
		}
		if form.Location != nil {
			event.Location = *form.Location
		}
		if form.Start != nil {
			event.Start = *form.Start
		}
		if form.Finish != nil {
			event.Finish = *form.Finish
		}
		if form.Public != nil {
			event.Public = *form.Public
		}
		if form.NotifyBefore != nil {
			event.NotifyBefore = *form.NotifyBefore
		}

		// Save the updated event back to the repository
		err = e.eventRepo.UpdateEvent(ctx, *event)
		if err != nil {
			log.Println("# Error updating event:", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		// Return a success response
		w.WriteHeader(http.StatusOK)
	}
}

// HandleDeleteEvent handles deleting an existing event.
func (e *EventHandler) HandleDeleteEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 3 || parts[2] == "" {
			log.Println("Invalid event ID in path")
			http.Error(w, "Invalid event ID in path", http.StatusBadRequest)
			return
		}

		idStr := parts[2] // Get the "1234" part
		eventID, err := uuid.Parse(idStr)
		if err != nil {
			log.Println("Invalid UUID")
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		// Retrieve the user ID from context
		user_id, ok := service.GetUserIDFromContext(r)
		if !ok {
			log.Println("Error: no user in context")
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		// Fetch the existing event from the repository
		ctx := r.Context()
		event, err := e.eventRepo.GetEventByID(ctx, eventID)
		if err != nil {
			log.Println("# Error fetching event:", err)
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}

		// Ensure the event belongs to the user (if applicable)
		if event.UserID != user_id {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Delete the event from the repository
		err = e.eventRepo.DeleteEvent(ctx, event.ID, user_id)
		if err != nil {
			log.Println("# Error deleting event:", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		// Return a success response
		w.WriteHeader(http.StatusNoContent) // No content since it's a delete operation
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
