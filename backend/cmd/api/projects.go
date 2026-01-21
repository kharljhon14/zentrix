package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kharljhon14/zentrix/internal/data"
	"github.com/kharljhon14/zentrix/internal/validator"
)

func (app application) createProjectHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CompanyID   string `json:"company_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		OwnerID     string `json:"owner_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	v := validator.New()

	v.ValidateUUID(input.CompanyID, "company_id")
	v.ValidateUUID(input.OwnerID, "owner_id")

	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	project := data.Project{
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		CompanyID:   uuid.MustParse(input.CompanyID),
		OwnerID:     uuid.MustParse(input.OwnerID),
	}

	project.Validate(v)
	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	app.models.Projects.Insert(&project)
}
