package handler

import (
	"golang-postgresql-auth-template/internal/app/repository"
	"golang-postgresql-auth-template/internal/app/service"
	"golang-postgresql-auth-template/internal/models"
	"log"
	"net/http"
	"strings"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo *repository.UserRepo
	appJwt   *service.AppJwt
}

func NewAuthHandler(userRepo *repository.UserRepo, appJwt *service.AppJwt) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		appJwt:   appJwt,
	}
}

func (a *AuthHandler) HandleSignup() http.HandlerFunc {
	type request struct {
		Email    string `json:"email" validate:"required,email,min=4,max=254"`
		Password string `json:"password" validate:"required,min=8,max=254"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		form, err := decodeJsonReq[request](r)
		if err != nil {
			log.Println("# Error decoding request", err)
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		err = service.Validate.Struct(form)
		if err != nil {
			log.Println("# Form validation error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		exists, err := a.userRepo.CheckIfUserExistsByEmail(ctx, form.Email)
		if err != nil {
			log.Println("# Error checking if user exists", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "A User with the associated email already exists", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("# Error hashing password", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		form.Password = string(hashedPassword)
		form.Email = strings.ToLower(form.Email)
		err = a.userRepo.CreateUser(ctx, form.Email, form.Password)
		if err != nil {
			log.Println("# Error creating user", err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (a *AuthHandler) HandleSignin() http.HandlerFunc {

	type request struct {
		Email    string `json:"email" validate:"required,min=4,max=254"`
		Password string `json:"password" validate:"required,min=8,max=254"`
	}
	type response struct {
		User    *models.User `json:"user"`
    Token   string       `json:"token"`
		Expires int64        `json:"exp"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		signInForm, err := decodeJsonReq[request](r)
		if err != nil {
			log.Println(err)
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}

		err = service.Validate.Struct(signInForm)
		if err != nil {
			log.Println(err)
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}
		signInForm.Email = strings.ToLower(signInForm.Email)
		user, err := a.userRepo.GetUserByEmail(r.Context(), signInForm.Email)
		if err != nil {
			log.Println(err)
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		if user == nil {
			log.Println("User does not exist")
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(signInForm.Password))
		if err != nil {
			log.Println(err)
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
			return
		}

		tokenString, exp, err := a.appJwt.CreateJwtWithClaims(user.ID)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		user.PasswordHash = ""
		signInResponse := response{
			User:    user,
			Token: tokenString,
			Expires: exp.Unix(),
		}

		a.appJwt.SetJwtCookie(w, tokenString, exp)
		sendJson(w, signInResponse, http.StatusOK)
	}
}

func (a *AuthHandler) HandleCheckIfSignedIn() http.HandlerFunc {
  type response struct {
    User      *models.User `json:"user"`
    ExpiresAt  int64       `json:"exp"`
  }
  return func(w http.ResponseWriter, r *http.Request) {
    claims := r.Context().Value(service.JwtClaimsContextKey).(*jwt.RegisteredClaims)
    user, err := a.userRepo.GetUserByID(r.Context(), claims.Subject)
    if err != nil {
      log.Println(err)
      http.Error(w, "Server Error", http.StatusInternalServerError)
      return
    }
    if user == nil {
      log.Println("User does not exist")
      http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
      return
    }
    user.PasswordHash = ""
    checkResponse := response{
      User:    user,
      ExpiresAt: claims.ExpiresAt.Unix(),
    }
	  sendJson(w, checkResponse, http.StatusOK)
  }
}

func (a *AuthHandler) HandleSignout(w http.ResponseWriter, r *http.Request) {
	a.appJwt.DeleteJwt(w)
	w.WriteHeader(http.StatusNoContent)
}
