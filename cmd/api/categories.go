package main

import (
	"TaskLogger/internal/data"
	"TaskLogger/internal/validator"
	"net/http"
)

func (bknd *backend) createCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name   string `json:"name"`
		Color  string `json:"color"`
		UserID int64  `json:"user_id"`
	}
	err := bknd.readJSON(w, r, &input)
	if err != nil {
		bknd.errBadRequest(w, r, err)
		return
	}
	vld := validator.NewValidator()
	category := &data.Categories{
		Name:   input.Name,
		Color:  input.Color,
		UserID: input.UserID,
	}
	if data.ValidateCategory(vld, category); !vld.Valid() {
		bknd.errFailedValidation(w, r, vld.Errors)
		return
	}
	_ = bknd.writeJSON(w, http.StatusCreated, wrapper{"input": input, "category": category}, nil)
}

func (bknd *backend) showCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bknd.readIdParam(r)
	if err != nil {
		bknd.errResourceNotFound(w, r)
		return
	}
	category := data.Categories{
		ID:     id,
		Name:   "Work",
		Color:  "#7DA925",
		UserID: 1,
	}
	err = bknd.writeJSON(w, http.StatusOK, wrapper{"category": category}, nil)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
	}
}
