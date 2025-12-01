package main

import "net/http"

func (app application) errorResponse(w http.ResponseWriter, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app application) serverErrorResponse(w http.ResponseWriter, err error) {
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, http.StatusInternalServerError, message)
}

func (app application) notFoundResponse(w http.ResponseWriter, message string) {
	app.errorResponse(w, http.StatusNotFound, message)
}
