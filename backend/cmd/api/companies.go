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

func (app application) createCompanyHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string  `json:"name"`
		Address string  `json:"address"`
		Email   string  `json:"email"`
		Image   *string `json:"image"`
		Website *string `json:"website"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	v := validator.New()

	company := &data.Company{
		Name:    input.Name,
		Address: input.Address,
		Email:   input.Email,
		Image:   input.Image,
		Website: input.Website,
	}

	v.Check(company.Name != "", "name", "name is required")
	v.Check(len(company.Name) <= 255, "name", "name must not exceed 255 characters")

	v.Check(company.Address != "", "address", "address is required")
	v.Check(len(company.Address) <= 255, "address", "address must not exceed 255 characters")

	v.Check(company.Email != "", "email", "email is required")

	if !validator.Matches(company.Email, validator.EmailRX) {
		v.AddError("email", "invalid email")
	}

	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	err = app.models.Companies.Insert(company)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "email already in use")
			app.failedValidationResponse(w, v.Errors)
		default:
			app.serverErrorResponse(w, err)
		}

		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/companies/%s", company.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"company": company}, headers)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app application) getCompanyByIDHandler(w http.ResponseWriter, r *http.Request) {
	IDParam := chi.URLParam(r, "id")

	v := validator.New()

	v.Check(IDParam != "", "id", "id is required")

	err := uuid.Validate(IDParam)
	if err != nil {
		v.AddError("id", "invalid id")
	}

	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	company, err := app.models.Companies.GetByID(uuid.MustParse(IDParam))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "company not found")
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"company": company}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app application) listCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 10, v)
	input.Sort = app.readString(qs, "sort", "id")
	input.SortSafeList =
		[]string{
			"id",
			"name",
			"address",
			"email",
			"created_at",
			"updated_at",
			"-id",
			"-name",
			"-address",
			"-created_at",
			"-updated_at",
		}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	companies, metadata, err := app.models.Companies.GetAll(input.Filters)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"companies": companies, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}

}
