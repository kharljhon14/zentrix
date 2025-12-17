package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/kharljhon14/zentrix/internal/data"
	"github.com/kharljhon14/zentrix/internal/validator"
)

func (app application) createContactHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		CompanyID string `json:"company_id"`
		Title     string `json:"title"`
		Status    string `json:"status"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	v := validator.New()

	contact := &data.Contact{
		Name:   input.Name,
		Email:  input.Email,
		Title:  input.Title,
		Status: input.Status,
	}

	v.ValidateUUID(input.CompanyID, "company_id")
	contact.ValidateContact(v)
	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	companyID := uuid.MustParse(input.CompanyID)
	contact.CompanyID = &companyID

	err = app.models.Contacts.Insert(contact)
	if err != nil {

		switch {
		case errors.Is(err, data.ErrInvalidUUID):
			v.AddError("company_id", "invalid ID")
			app.failedValidationResponse(w, v.Errors)
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "email already in use")
			app.failedValidationResponse(w, v.Errors)
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/contacts/%s", contact.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"data": contact}, headers)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
