package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kharljhon14/zentrix/internal/validator"
)

func (app application) getProductsByQuoteIDHandler(w http.ResponseWriter, r *http.Request) {
	IDParam := chi.URLParam(r, "id")

	v := validator.New()
	v.ValidateUUID(IDParam, "id")

	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	products, err := app.models.Products.GetProductsByQuoteID(uuid.MustParse(IDParam))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "no products found with given quote ID")
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"data": products}, nil)
}
