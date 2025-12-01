package main

import (
	"net/http"

	"github.com/kharljhon14/zentrix/internal/data"
	"github.com/kharljhon14/zentrix/internal/validator"
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

	v := validator.New()
	v.Check(input.Name != "", "name", "name is required")
	v.Check(len(input.Name) >= 1, "name", "name is required")
	v.Check(len(input.Name) <= 255, "name", "name must not exceed 255 character")
	v.Check(input.Plan != "", "plan", "plan is required")
	v.Check(len(input.Plan) >= 1, "plan", "plan is required")
	v.Check(len(input.Plan) <= 255, "plan", "plan must not exceed 255 character")

	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
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
