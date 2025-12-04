package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kharljhon14/zentrix/internal/data"
	"github.com/kharljhon14/zentrix/internal/validator"
)

func (app application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Password  string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	user := &data.User{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Role:      "admin",
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		fmt.Println(err)
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "email already exists")
			app.failedValidationResponse(w, v.Errors)
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user, "token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		PlainTextToken string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	v := validator.New()

	if data.ValidatePlainTextToken(v, input.PlainTextToken); !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	user, err := app.models.Tokens.GetForToken(data.ScopeActivation, input.PlainTextToken)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, v.Errors)
		default:
			app.serverErrorResponse(w, err)
		}

		return
	}

	user.Activated = true

	err = app.models.Users.Update(user)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}

}
