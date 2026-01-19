package main

import (
	"database/sql"
	"errors"
	"fmt"
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

// Title     string    `json:"title"`
// 	UnitPrice int       `json:"unit_price"`
// 	Quantity  int       `json:"quantity"`
// 	Discount  int       `json:"discount"`

func (app application) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	IDParam := chi.URLParam(r, "id")

	var input struct {
		Title     *string `json:"title"`
		UnitPrice *int    `json:"unit_price"`
		Quantity  *int    `json:"quantity"`
		Discount  *int    `json:"discount"`
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
	v.ValidateUUID(IDParam, "id")

	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	product, err := app.models.Products.GetProductByID(uuid.MustParse(IDParam))
	if err != nil {
		fmt.Println(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "product not found")
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	if input.Title != nil {
		product.Title = *input.Title
	}
	if input.UnitPrice != nil {
		product.UnitPrice = *input.UnitPrice
	}
	if input.Quantity != nil {
		product.Quantity = *input.Quantity
	}
	if input.Discount != nil {
		product.Discount = *input.Discount
	}

	if product.ValidateProduct(v); !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	product, err = app.models.Products.Update(product)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "product not found")
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": product}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}

}

func (app application) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	IDParam := chi.URLParam(r, "id")

	v := validator.New()

	v.ValidateUUID(IDParam, "id")

	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	err := app.models.Products.Delete(uuid.MustParse(IDParam))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "product not found")
		default:
			app.serverErrorResponse(w, err)
		}
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "product deleted successfully"}, nil)
}
