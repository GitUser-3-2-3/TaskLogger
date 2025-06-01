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
		UserID string `json:"user_id"`
	}
	err := bknd.readJSON(w, r, &input)
	if err != nil {
		bknd.errBadRequest(w, r, err)
		return
	}
	vld := validator.NewValidator()

	ctg := &data.Categories{Name: input.Name,
		Color:  input.Color,
		UserID: input.UserID,
	}
	if data.ValidateCategory(vld, ctg); !vld.Valid() {
		bknd.errFailedValidation(w, r, vld.Errors)
		return
	}
	err = bknd.models.Category.Insert(ctg)
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
	headers.Set("Location", fmt.Sprintf("/v1/category/%d", ctg.ID))

	err = bknd.writeJSON(w, http.StatusCreated, wrapper{"categoryId": ctg.ID}, headers)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
	}
}

func (bknd *backend) showCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bknd.readIdParam(r, "categoryId")
	if err != nil {
		bknd.errResourceNotFound(w, r)
		return
	}
	ctg, err := bknd.models.Category.GetByCtgId(id)
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

// todo -> add user id retrieval from context instead of url.

func (bknd *backend) showCategoriesForUserIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bknd.readUUIDParam(r, "userId")
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

// todo -> add user id retrieval from context instead of url.
// Inefficient method!

func (bknd *backend) showCategoryDetailsForUserIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bknd.readUUIDParam(r, "userId")
	if err != nil {
		bknd.errResourceNotFound(w, r)
		return
	}
	ctgList, err := bknd.models.Category.GetAllByUserId(id)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
		return
	}
	var taskList []*data.Tasks
	var ctgDetails []wrapper

	for _, ctg := range ctgList {
		taskList, err = bknd.models.Tasks.GetAllByCategory(ctg.ID)
		if err != nil {
			bknd.errInternalServerError(w, r, err)
			return
		}
		ctgDetails = append(ctgDetails, wrapper{ctg.Name: taskList})
	}
	err = bknd.writeJSON(w, http.StatusOK, wrapper{"categories": ctgDetails}, nil)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
	}
}

func (bknd *backend) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bknd.readIdParam(r, "categoryId")
	if err != nil {
		bknd.errResourceNotFound(w, r)
		return
	}
	category, err := bknd.models.Category.GetByCtgId(id)
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
		UserID *string `json:"-"`
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
	id, err := bknd.readIdParam(r, "categoryId")
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
