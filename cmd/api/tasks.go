package main

import (
	"TaskLogger/internal/data"
	"TaskLogger/internal/validator"
	"errors"
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
		UserID      string          `json:"user_id"`
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
		ImageUrl:    input.Image,
		Deadline:    input.Deadline,
		UserID:      input.UserID,
		CategoryID:  input.CategoryID,
	}
	if data.ValidateTask(vld, task); !vld.Valid() {
		bknd.errFailedValidation(w, r, vld.Errors)
		return
	}
	err = bknd.models.Tasks.Insert(task)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/tasks/%s", task.ID))

	err = bknd.writeJSON(w, http.StatusCreated, wrapper{"id": task.ID}, headers)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
	}
}

func (bknd *backend) showTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bknd.readUUIDParam(r, "id")
	if err != nil {
		bknd.errResourceNotFound(w, r)
		return
	}
	task, err := bknd.models.Tasks.GetByTaskId(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			bknd.errResourceNotFound(w, r)
		default:
			bknd.errInternalServerError(w, r, err)
		}
		return
	}
	err = bknd.writeJSON(w, http.StatusOK, wrapper{"task": task}, nil)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
	}
}

func (bknd *backend) showTasksByCategory(w http.ResponseWriter, r *http.Request) {
	id, err := bknd.readIdParam(r, "id")
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
	tasks, err := bknd.models.Tasks.GetAllByCategory(id)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
		return
	}
	err = bknd.writeJSON(w, http.StatusOK, wrapper{"category": ctg.Name, "tasks": tasks}, nil)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
	}
}

func (bknd *backend) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bknd.readUUIDParam(r, "id")
	if err != nil {
		bknd.errResourceNotFound(w, r)
		return
	}
	task, err := bknd.models.Tasks.GetByTaskId(id)
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
		Name        *string          `json:"name"`
		Description *string          `json:"description"`
		Status      *data.StatusType `json:"status"`
		Image       *string          `json:"image"`
		Priority    *int             `json:"priority"`
		Deadline    *time.Time       `json:"deadline"`
		UserID      *string          `json:"-"`
		CategoryID  *int64           `json:"category_id"`
	}
	err = bknd.readJSON(w, r, &input)
	if err != nil {
		bknd.errBadRequest(w, r, err)
		return
	}
	task.ApplyPartialUpdatesToTask(input.Name, input.Description,
		input.Image, input.Status, input.Priority, input.Deadline, input.UserID, input.CategoryID)

	vld := validator.NewValidator()
	if data.ValidateTask(vld, task); !vld.Valid() {
		bknd.errFailedValidation(w, r, vld.Errors)
		return
	}
	err = bknd.models.Tasks.Update(task)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			bknd.errResourceNotFound(w, r)
		default:
			bknd.errInternalServerError(w, r, err)
		}
		return
	}
	err = bknd.writeJSON(w, http.StatusOK, wrapper{"task": task}, nil)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
	}
}

func (bknd *backend) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bknd.readUUIDParam(r, "id")
	if err != nil {
		bknd.errResourceNotFound(w, r)
		return
	}
	err = bknd.models.Tasks.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			bknd.errResourceNotFound(w, r)
		default:
			bknd.errInternalServerError(w, r, err)
		}
		return
	}
	err = bknd.writeJSON(w, http.StatusOK, wrapper{"message": "task deleted"}, nil)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
	}
}
