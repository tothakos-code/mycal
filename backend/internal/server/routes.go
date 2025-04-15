package server

import (
	"golang-postgresql-auth-template/internal/app"
	"net/http"
)

func RegisterRoutes(app *app.Application) http.Handler {
	router := http.NewServeMux()
	middlewares := CreateMiddlewareStack(
		Cors,
		Timeout,
		Logger,
	)

	// Public Routes
	router.HandleFunc("POST /auth/signup", app.AuthHandler.HandleSignup())
	router.HandleFunc("POST /auth/signin", app.AuthHandler.HandleSignin())
	router.HandleFunc("POST /auth/signout", app.AuthHandler.HandleSignout)
	router.HandleFunc("POST /auth/me", app.AppJwt.IsAuthenticatedMiddleware(app.AuthHandler.HandleCheckIfSignedIn()))

	v1 := http.NewServeMux()
	// V1 Not protected routes
	v1.HandleFunc("GET /event/last-30", app.EventHandler.HandleListPublicEvents())

	// V1 Protected Routes
	// Events
	v1.HandleFunc("GET /event/last-30-private", app.AppJwt.IsAuthenticatedMiddleware(app.EventHandler.HandleListPrivateEvents()))
	v1.HandleFunc("POST /event", app.AppJwt.IsAuthenticatedMiddleware(app.EventHandler.HandleCreateEvent()))
	v1.HandleFunc("PUT /event/{id}", app.AppJwt.IsAuthenticatedMiddleware(app.EventHandler.HandleUpdateEvent()))
	v1.HandleFunc("DELETE /event/{id}", app.AppJwt.IsAuthenticatedMiddleware(app.EventHandler.HandleDeleteEvent()))

	// Users
	v1.HandleFunc("PUT /user/{id}", app.AppJwt.IsAuthenticatedMiddleware(app.UserHandler.HandleUpdateUser()))
	v1.HandleFunc("GET /user/{id}", app.UserHandler.HandleFetchUser())
	// v1.HandleFunc("GET /invitation/pending")

	// Add middleware to protected routes
	router.Handle("/v1/", http.StripPrefix("/v1", v1))
	return middlewares(router)
}
