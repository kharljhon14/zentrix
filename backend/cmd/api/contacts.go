package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
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

func (app application) getContactByIDHandler(w http.ResponseWriter, r *http.Request) {
	IDParam := chi.URLParam(r, "id")

	v := validator.New()

	v.Check(IDParam != "", "id", "id is required")
	v.ValidateUUID(IDParam, "id")

	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	contact, err := app.models.Contacts.GetByIDWithCompanyName(uuid.MustParse(IDParam))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "contact not found")
		default:
			app.serverErrorResponse(w, err)
		}

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": contact}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app application) listContactsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CompanyID *uuid.UUID
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 10, v)
	input.Sort = app.readString(qs, "sort", "-created_at")
	input.SortSafeList =
		[]string{
			"id",
			"name",
			"email",
			"company_name",
			"title",
			"status",
			"created_at",
			"updated_at",
			"-id",
			"-name",
			"-email",
			"-company_name",
			"-title",
			"-status",
			"-created_at",
			"-updated_at",
		}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	if qs.Has("company_id") {
		companyID, err := uuid.Parse(qs.Get("company_id"))
		if err != nil {
			input.CompanyID = nil
		} else {
			input.CompanyID = &companyID
		}
	}

	contacts, metadata, err := app.models.Contacts.GetAll(input.Filters, input.CompanyID)
	if err != nil {
		fmt.Println(err)

		app.serverErrorResponse(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": contacts, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
