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

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
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

	companyID := uuid.MustParse(input.CompanyID)
	_, err = app.models.Companies.GetByID(companyID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "company not found")
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	preparedBy := uuid.MustParse(input.PreparedBy)
	_, err = app.models.Users.GetByID(preparedBy)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "prepared By ID not found")
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	prepareFor := uuid.MustParse(input.PreparedFor)
	_, err = app.models.Contacts.GetByID(prepareFor)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "prepared for ID not found")
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	quote.CompanyID = companyID
	quote.PreparedBy = preparedBy
	quote.PreparedFor = prepareFor

	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	err = app.models.Quotes.Insert(&quote)
	if err != nil {

		fmt.Println(err)
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

	for i, product := range products {
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
		products[i] = product
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{
		"quote":    quote,
		"products": products,
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app application) getQuoteByIDHandler(w http.ResponseWriter, r *http.Request) {
	IDParam := chi.URLParam(r, "id")

	v := validator.New()
	v.Check(IDParam != "", "id", "id is required")
	if v.ValidateUUID(IDParam, "id"); !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	quote, err := app.models.Quotes.GetByID(uuid.MustParse(IDParam))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "quote not found")
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	products, err := app.models.Products.GetProductsByQuoteID(quote.ID)
	if err != nil {
		fmt.Println(err)
		app.serverErrorResponse(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"data": envelope{
		"quote":    quote,
		"products": products,
	}}, nil)

}

func (app application) listQuotesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 10, v)
	input.Sort = app.readString(qs, "sort", "-created_at")
	input.SortSafeList = []string{
		"id",
		"company_id",
		"name",
		"prepared_by",
		"prepared_for",
		"stage",
		"created_at",
		"updated_at",
		"-id",
		"-company_id",
		"-name",
		"-prepared_by",
		"-prepared_for",
		"-stage",
		"-created_at",
		"-updated_at",
	}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	quotes, metadata, err := app.models.Quotes.GetAll(input.Filters)
	if err != nil {
		fmt.Println(err)
		app.serverErrorResponse(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": quotes, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app application) updateQuoteHandler(w http.ResponseWriter, r *http.Request) {
	IDParam := chi.URLParam(r, "id")

	var input struct {
		Name        *string `json:"name"`
		CompanyID   *string `json:"company_id"`
		Stage       *string `json:"stage"`
		Notes       *string `json:"notes"`
		PreparedBy  *string `json:"prepared_by"`
		PreparedFor *string `json:"prepared_for"`
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
	if input.CompanyID != nil {
		v.ValidateUUID(*input.CompanyID, "company_id")
	}
	if input.PreparedBy != nil {
		v.ValidateUUID(*input.PreparedBy, "prepared_by")
	}
	if input.PreparedFor != nil {
		v.ValidateUUID(*input.PreparedFor, "prepared_for")
	}

	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	quote, err := app.models.Quotes.GetByID(uuid.MustParse(IDParam))
	if err != nil {
		fmt.Println(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "quote not found")
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	if input.Name != nil {
		quote.Name = *input.Name
	}
	if input.CompanyID != nil {

		quote.CompanyID = uuid.MustParse(*input.CompanyID)

		_, err = app.models.Companies.GetByID(quote.CompanyID)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				app.notFoundResponse(w, "company not found")
			default:
				app.serverErrorResponse(w, err)
			}
			return
		}
	}
	if input.Stage != nil {
		quote.Stage = *input.Stage
	}
	if input.Notes != nil {
		quote.Notes = *input.Notes
	}
	if input.PreparedBy != nil {
		quote.PreparedBy = uuid.MustParse(*input.PreparedBy)

		_, err = app.models.Users.GetByID(quote.PreparedBy)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				app.notFoundResponse(w, "user not found")
			default:
				app.serverErrorResponse(w, err)
			}
			return
		}
	}
	if input.PreparedFor != nil {
		quote.PreparedFor = uuid.MustParse(*input.PreparedFor)

		_, err = app.models.Contacts.GetByID(quote.PreparedFor)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				app.notFoundResponse(w, "contact not found")
			default:
				app.serverErrorResponse(w, err)
			}
			return
		}
	}

	quote.ValidateQuote(v)
	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	err = app.models.Quotes.Update(quote)
	if err != nil {
		fmt.Println(err)
		app.serverErrorResponse(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": quote}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app application) deleteQuoteHandler(w http.ResponseWriter, r *http.Request) {
	IDParam := chi.URLParam(r, "id")

	v := validator.New()
	v.ValidateUUID(IDParam, "id")
	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	err := app.models.Quotes.Delete(uuid.MustParse(IDParam))
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.notFoundResponse(w, "quote not found")
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "quote deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
