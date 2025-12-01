package main

import "net/http"

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	env := envelope{
		"status":  "available",
		"version": version,
		"env":     app.config.env,
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		w.Write([]byte("error"))
	}
}
