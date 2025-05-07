package main

import (
	"TaskLogger/internal/data"
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
	_ = bknd.writeJSON(w, http.StatusCreated, wrapper{"category": input}, nil)
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
