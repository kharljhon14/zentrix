package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/kharljhon14/zentrix/internal/data"
	"github.com/kharljhon14/zentrix/internal/validator"
)

func (app application) createQuoteHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		CompanyID   string `json:"company_id"`
		SalesTax    int    `json:"sales_tax"`
		Stage       string `json:"stage"`
		Notes       string `json:"notes"`
		PreparedBy  string `json:"prepared_by"`
		PreparedFor string `json:"prepared_for"`
		Products    []struct {
			Title     string `json:"title"`
			UnitPrice int    `json:"unit_price"`
			Quantity  int    `json:"quantity"`
			Discount  int    `json:"discount"`
		} `json:"products"`
	}

	v := validator.New()

	v.ValidateUUID(input.CompanyID, "company_id")
	v.ValidateUUID(input.PreparedBy, "prepared_by")
	v.ValidateUUID(input.PreparedFor, "prepared_for")

	quote := data.Quote{
		Name:     input.Name,
		SalesTax: input.SalesTax,
		Stage:    input.Stage,
		Notes:    input.Notes,
	}

	quote.ValidateQuote(v)

	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	companyID := uuid.MustParse(input.CompanyID)
	_, err := app.models.Companies.GetByID(companyID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "company not found")
		default:
			app.serverErrorResponse(w, err)
			return
		}
	}

	preparedBy := uuid.MustParse(input.PreparedBy)
	_, err = app.models.Users.GetByID(preparedBy)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "prepared By ID not found")
		default:
			app.serverErrorResponse(w, err)
			return
		}
	}

	prepareFor := uuid.MustParse(input.PreparedFor)
	_, err = app.models.Contacts.GetByID(prepareFor)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "prepared for ID not found")
		default:
			app.serverErrorResponse(w, err)
			return
		}
	}

	quote.CompanyID = companyID
	quote.PreparedBy = preparedBy
	quote.PreparedFor = prepareFor

	err = app.models.Quotes.Insert(&quote)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	var products []data.Product
	for _, productInput := range input.Products {
		product := data.Product{
			QuoteID:   quote.ID,
			Title:     productInput.Title,
			UnitPrice: productInput.UnitPrice,
			Quantity:  productInput.Quantity,
			Discount:  productInput.Discount,
		}
		products = append(products, product)
	}

	for _, product := range products {
		product.ValidateProduct(v)

		if !v.Valid() {
			app.failedValidationResponse(w, v.Errors)
			//TODO Remove the inserted Quote
			return
		}

		err := app.models.Products.Insert(&product)
		if err != nil {
			app.serverErrorResponse(w, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{
		"quote":    quote,
		"products": products,
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
