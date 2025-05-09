package main

import (
	"TaskLogger/internal/data"
	"TaskLogger/internal/validator"
	"errors"
	"fmt"
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
	id, err := bknd.models.Category.Insert(category)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEntry):
			bknd.errDuplicateEntryFound(w, r, err)
		default:
			bknd.errInternalServerError(w, r, err)
		}
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/category/%d", id))

	err = bknd.writeJSON(w, http.StatusCreated, wrapper{"id": id}, headers)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
	}
}

func (bknd *backend) showCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bknd.readIdParam(r, "id")
	if err != nil {
		bknd.errResourceNotFound(w, r)
		return
	}
	ctg, err := bknd.models.Category.GetById(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			bknd.errResourceNotFound(w, r)
		default:
			bknd.errInternalServerError(w, r, err)
		}
		return
	}
	err = bknd.writeJSON(w, http.StatusOK, wrapper{"category": ctg}, nil)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
	}
}

func (bknd *backend) showCategoriesByUserIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bknd.readIdParam(r, "userId")
	if err != nil {
		bknd.errResourceNotFound(w, r)
		return
	}
	ctgList, err := bknd.models.Category.GetAllByUserId(id)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
		return
	}
	err = bknd.writeJSON(w, http.StatusOK, wrapper{"categories": ctgList}, nil)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
	}
}

func (bknd *backend) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bknd.readIdParam(r, "id")
	if err != nil {
		bknd.errResourceNotFound(w, r)
		return
	}
	category, err := bknd.models.Category.GetById(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			bknd.errResourceNotFound(w, r)
		default:
			bknd.errInternalServerError(w, r, err)
		}
		return
	}
	var input struct {
		Name   *string `json:"name"`
		Color  *string `json:"color"`
		UserID *int64  `json:"-"`
	}
	err = bknd.readJSON(w, r, &input)
	if err != nil {
		bknd.errBadRequest(w, r, err)
		return
	}
	category.ApplyPartialUpdatesToCtg(input.Name, input.Color, input.UserID)

	vld := validator.NewValidator()
	if data.ValidateCategory(vld, category); !vld.Valid() {
		bknd.errFailedValidation(w, r, vld.Errors)
		return
	}
	err = bknd.models.Category.Update(category)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			bknd.errResourceNotFound(w, r)
		default:
			bknd.errInternalServerError(w, r, err)
		}
		return
	}
	err = bknd.writeJSON(w, http.StatusOK, wrapper{"category": category}, nil)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
	}
}

func (bknd *backend) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bknd.readIdParam(r, "id")
	if err != nil {
		bknd.errResourceNotFound(w, r)
		return
	}
	err = bknd.models.Category.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			bknd.errResourceNotFound(w, r)
		default:
			bknd.errInternalServerError(w, r, err)
		}
		return
	}
	err = bknd.writeJSON(w, http.StatusOK, wrapper{"message": "deletion successful"}, nil)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
	}
}
