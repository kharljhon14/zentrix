package main

import (
	"errors"
	"fmt"
	"net/http"

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

	err = app.writeJSON(w, http.StatusCreated, envelope{"company": company}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
