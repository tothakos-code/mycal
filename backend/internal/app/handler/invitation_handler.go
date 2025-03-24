package handler

import (
	"golang-postgresql-auth-template/internal/app/repository"
	"log"
	"net/http"
	"encoding/json"
)

type InvitationHandler struct {
	invitationRepo *repository.InvitationRepo
}

func NewInvitationHandler(invitationRepo *repository.InvitationRepo) *InvitationHandler {
	return &InvitationHandler{invitationRepo: invitationRepo}
}

func (i *InvitationHandler) HandleCreateInvitation() http.HandlerFunc {
	type request struct {
		UserID  string `json:"user_id" validate:"required,uuid"`
		EventID string `json:"event_id" validate:"required,uuid"`
		Status  string `json:"status" validate:"required,oneof=pending accepted declined"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
    form, err := decodeJsonReq[request](r)
		if err != nil {
			log.Println("# Error decoding request", err)
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		err = i.invitationRepo.CreateInvitation(ctx, form.UserID, form.EventID, form.Status)
		if err != nil {
			log.Println("# Error creating invitation:", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (i *InvitationHandler) HandleListInvitedUsersByEventID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		eventID := r.URL.Query().Get("event_id")
		if eventID == "" {
			http.Error(w, "Missing event_id parameter", http.StatusBadRequest)
			return
		}

		users, err := i.invitationRepo.ListInvitedUsersByEventID(ctx, eventID)
		if err != nil {
			log.Println("# Error listing invited users:", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}
