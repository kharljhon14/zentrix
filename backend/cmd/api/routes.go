package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/healthcheck", app.healthCheckHandler)

	// User auth
	r.Post("/register", app.registerUserHandler)
	r.Put("/activate", app.activateUserHandler)

	return r
}
