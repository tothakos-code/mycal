package handler

import (
	"github.com/google/uuid"
	"golang-postgresql-auth-template/internal/app/repository"
	"golang-postgresql-auth-template/internal/models"
	"log"
	"net/http"
	"time"
)

type NotificationHandler struct {
	notificationRepo *repository.NotificationRepo
}

func NewNotificationHandler(notificationRepo *repository.NotificationRepo) *NotificationHandler {
	return &NotificationHandler{notificationRepo: notificationRepo}
}

func (n *NotificationHandler) HandleCreateNotification() http.HandlerFunc {
	type request struct {
		UserID  uuid.UUID `json:"user_id" validate:"required,uuid"`
		EventID uuid.UUID `json:"event_id" validate:"required,uuid"`
		Title   string `json:"title" validate:"required"`
		Message string `json:"message" validate:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
    form, err := decodeJsonReq[request](r)
    if err != nil {
      log.Println("# Error decoding request", err)
      http.Error(w, "Invalid Request", http.StatusBadRequest)
      return
    }

		notification := models.Notification{
			UserID:    form.UserID,
			EventID:   form.EventID,
			Title:     form.Title,
			Message:   form.Message,
			ShownAt:   time.Now(),
			CreatedAt: time.Now(),
		}

		err = n.notificationRepo.CreateNotification(ctx, notification)
		if err != nil {
			log.Println("# Error creating notification:", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
