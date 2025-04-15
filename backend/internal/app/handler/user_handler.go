package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"golang-postgresql-auth-template/internal/app/repository"
	"golang-postgresql-auth-template/internal/app/service"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

type UserHandler struct {
	userRepo *repository.UserRepo
}

func NewUserHandler(userRepo *repository.UserRepo) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// HandleEditEvent handles updating an existing event.
func (u *UserHandler) HandleUpdateUser() http.HandlerFunc {
	type request struct {
		FirstName *string `json:"firstname"`
		SurName   *string `json:"surname"`
		Password  *string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		log.Println(parts)
		if len(parts) != 3 || parts[2] == "" {
			log.Println("Invalid user ID in path")
			http.Error(w, "Invalid user ID in path", http.StatusBadRequest)
			return
		}

		idStr := parts[2]
		userID, err := uuid.Parse(idStr)
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
		requester_user_id, ok := service.GetUserIDFromContext(r)
		if !ok {
			log.Println("Error: no user in context", err)
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		// Fetch the existing event from the repository
		user, err := u.userRepo.GetUserByID(ctx, userID.String())
		if err != nil {
			log.Println("# Error fetching event:", err)
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}

		// Ensure the event belongs to the user (if applicable)
		if user.ID != requester_user_id {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Update fields in the event if they are provided in the request
		if form.FirstName != nil {
			user.FirstName = *form.FirstName
		}
		if form.SurName != nil {
			user.SurName = *form.SurName
		}
		if form.Password != nil {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*form.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Println("# Error hashing password", err)
				http.Error(w, "Server Error", http.StatusInternalServerError)
				return
			}
			user.PasswordHash = string(hashedPassword)
		}

		// Save the updated event back to the repository
		err = u.userRepo.UpdateUser(ctx, *user)
		if err != nil {
			log.Println("# Error updating user:", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		// Return a success response
		w.WriteHeader(http.StatusOK)
	}
}
func (u *UserHandler) HandleFetchUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		log.Println(parts)
		if len(parts) != 3 || parts[2] == "" {
			log.Println("Invalid user ID in path")
			http.Error(w, "Invalid user ID in path", http.StatusBadRequest)
			return
		}

		idStr := parts[2]
		_, err := uuid.Parse(idStr)
		if err != nil {
			log.Println("Invalid UUID")
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		user, err := u.userRepo.GetUserByID(ctx, idStr)
		if err != nil {
			log.Println("Error getting user:", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)

	}
}
