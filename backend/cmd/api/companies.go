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
		Name         string  `json:"name"`
		Address      string  `json:"address"`
		SalesOwner   string  `json:"sales_owner"`
		Email        string  `json:"email"`
		CompanySize  string  `json:"company_size"`
		Industry     string  `json:"industry"`
		BusinessType string  `json:"business_type"`
		Country      string  `json:"country"`
		Image        *string `json:"image"`
		Website      *string `json:"website"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	v := validator.New()

	company := &data.Company{
		Name:         input.Name,
		Address:      input.Address,
		Email:        input.Email,
		CompanySize:  input.CompanySize,
		Industry:     input.Industry,
		BusinessType: input.BusinessType,
		Country:      input.Country,
		Image:        input.Image,
		Website:      input.Website,
	}

	// Validate the input values including the sales owner id format
	v.ValidateUUID(input.SalesOwner, "sale_owner")
	if data.ValidateCompany(v, company); !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	company.SalesOwner = uuid.MustParse(input.SalesOwner)

	err = app.models.Companies.Insert(company)
	if err != nil {

		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "email already in use")
			app.failedValidationResponse(w, v.Errors)
		default:
			//TODO check for existing sales_owner
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
	v.ValidateUUID(IDParam, "id")

	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	company, err := app.models.Companies.GetByIDWithSalesOwner(uuid.MustParse(IDParam))
	if err != nil {
		fmt.Println(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "company not found")
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": company}, nil)
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
	input.Sort = app.readString(qs, "sort", "-created_at")
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

	err = app.writeJSON(w, http.StatusOK, envelope{"data": companies, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}

}

func (app application) updatedCompanyHandler(w http.ResponseWriter, r *http.Request) {
	IDParam := chi.URLParam(r, "id")

	var input struct {
		Name         *string `json:"name"`
		Address      *string `json:"address"`
		SalesOwner   *string `json:"sales_owner"`
		Email        *string `json:"email"`
		CompanySize  *string `json:"company_size"`
		Industry     *string `json:"industry"`
		BusinessType *string `json:"business_type"`
		Country      *string `json:"country"`
		Image        *string `json:"image"`
		Website      *string `json:"website"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	if app.isAllNil(input) {
		app.badRequestResponse(w, errors.New("body must not be empty"))
		return
	}

	v := validator.New()

	// Validate UUIDs
	v.Check(IDParam != "", "id", "id is required")
	v.ValidateUUID(IDParam, "id")
	if input.SalesOwner != nil {
		v.ValidateUUID(*input.SalesOwner, "sales_owner")
	}

	company, err := app.models.Companies.GetByID(uuid.MustParse(IDParam))
	if err != nil {
		fmt.Println(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "company not found")
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	// Merge only the non nil fields from the input
	// into the existing company record.

	if input.Name != nil {
		company.Name = *input.Name
	}

	if input.Address != nil {
		company.Address = *input.Address
	}

	if input.Email != nil {
		company.Email = *input.Email
	}

	if input.CompanySize != nil {
		company.CompanySize = *input.CompanySize
	}

	if input.BusinessType != nil {
		company.BusinessType = *input.BusinessType
	}

	if input.Industry != nil {
		company.Industry = *input.Industry
	}

	if input.Country != nil {
		company.Country = *input.Country
	}

	if input.Image != nil {
		company.Image = input.Image
	}

	if input.Website != nil {
		company.Website = input.Website
	}

	if data.ValidateCompany(v, company); !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	if input.SalesOwner != nil {
		salesOwnerID := uuid.MustParse(*input.SalesOwner)
		_, err := app.models.Users.GetByID(salesOwnerID)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				app.notFoundResponse(w, "sales_owner not found")
			default:
				app.serverErrorResponse(w, err)
			}

			return
		}

		company.SalesOwner = salesOwnerID
	}

	err = app.models.Companies.Update(company)
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

	err = app.writeJSON(w, http.StatusOK, envelope{"data": company}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app application) deleteCompanyHandler(w http.ResponseWriter, r *http.Request) {
	IDParam := chi.URLParam(r, "id")

	v := validator.New()

	v.Check(IDParam != "", "id", "id is required")
	v.ValidateUUID(IDParam, "id")

	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	err := app.models.Companies.Delete(uuid.MustParse(IDParam))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "company not found")
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "company deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
