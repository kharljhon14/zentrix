package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(httprate.LimitByIP(100, time.Minute))

	r.Get("/healthcheck", app.healthCheckHandler)

	// User auth
	r.Post("/register", app.registerUserHandler)
	r.Put("/activate", app.activateUserHandler)

	return r
}
