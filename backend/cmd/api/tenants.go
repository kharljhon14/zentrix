package main

import (
	"net/http"

	"github.com/kharljhon14/zentrix/internal/data"
)

func (app application) registerTenantHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
		Plan string `json:"plan"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	tenant := &data.Tenant{
		Name: input.Name,
		Plan: input.Plan,
	}

	err = app.models.Tenants.Insert(tenant)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"tenant": tenant}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}

}
