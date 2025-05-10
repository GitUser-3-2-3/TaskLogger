package main

import (
	"TaskLogger/internal/data"
	"TaskLogger/internal/validator"
	"fmt"
	"net/http"
	"time"
)

func (bknd *backend) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Status      data.StatusType `json:"status"`
		Priority    int             `json:"priority"`
		Image       string          `json:"image"`
		Deadline    *time.Time      `json:"deadline"`
		UserID      int64           `json:"user_id"`
		CategoryID  *int64          `json:"category_id"`
	}
	err := bknd.readJSON(w, r, &input)
	if err != nil {
		bknd.errBadRequest(w, r, err)
		return
	}
	vld := validator.NewValidator()

	task := &data.Tasks{
		Name:        input.Name,
		Description: input.Description,
		Status:      input.Status,
		Priority:    input.Priority,
		Image:       input.Image,
		Deadline:    input.Deadline,
		UserID:      input.UserID,
		CategoryID:  input.CategoryID,
	}
	if data.ValidateTask(vld, task); !vld.Valid() {
		bknd.errFailedValidation(w, r, vld.Errors)
		return
	}
	id, err := bknd.models.Tasks.Insert(task)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/tasks/%d", id))

	err = bknd.writeJSON(w, http.StatusCreated, wrapper{"id": id}, headers)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
	}
}

func (bknd *backend) showTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bknd.readIdParam(r, "id")
	if err != nil {
		bknd.errResourceNotFound(w, r)
		return
	}
	deadline := time.Now().Add(24 * time.Hour)
	task := data.Tasks{
		ID:            id,
		Name:          "Project",
		Description:   "",
		Status:        "",
		Priority:      1,
		Image:         "",
		TotalDuration: 0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Deadline:      &deadline,
		UserID:        1,
		CategoryID:    nil,
	}
	err = bknd.writeJSON(w, http.StatusOK, wrapper{"task": task}, nil)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
	}
}
