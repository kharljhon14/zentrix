package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(httprate.LimitByIP(100, time.Minute))

	r.Get("/healthcheck", app.healthCheckHandler)

	// User auth
	r.Post("/register", app.registerUserHandler)
	r.Put("/activate", app.activateUserHandler)

	// Companies
	r.Post("/companies", app.createCompanyHandler)
	r.Get("/companies", app.listCompaniesHandler)
	r.Get("/companies/{id}", app.getCompanyByIDHandler)
	r.Patch("/companies/{id}", app.updatedCompanyHandler)
	r.Delete("/companies/{id}", app.deleteCompanyHandler)
	return r
}
