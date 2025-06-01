package main

import (
	"TaskLogger/internal/data"
	"TaskLogger/internal/validator"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// todo -> Fetch taskID from uri, unless there is a better approach.
// 	     Add validation for fields as well as checking whether a task exists or not.

func (bknd *backend) createSessionHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TaskID    string    `json:"task_id"`
		StartedAt time.Time `json:"started_at"`
		EndedAt   time.Time `json:"ended_at"`
		Note      string    `json:"note"`
	}
	err := bknd.readJSON(w, r, &input)
	if err != nil {
		bknd.errBadRequest(w, r, err)
		return
	}
	if input.EndedAt.Before(input.StartedAt) {
		bknd.errBadRequest(w, r, fmt.Errorf("started-at must be before ended-at"))
		return
	}
	duration := int(input.EndedAt.Sub(input.StartedAt).Minutes())

	session := &data.Session{
		TaskID:    input.TaskID,
		StartedAt: input.StartedAt,
		EndedAt:   input.EndedAt,
		Duration:  duration,
		Note:      input.Note,
	}
	vld := validator.NewValidator()

	if data.ValidateSession(vld, session); !vld.Valid() {
		bknd.errFailedValidation(w, r, vld.Errors)
		return
	}
	err = bknd.withTransaction(func(tx *sql.Tx) error {
		var txErr error
		task, err := bknd.models.Tasks.GetByTaskIdTx(tx, session.TaskID)
		if err != nil {
			return err
		}
		txErr = bknd.models.Session.InsertTx(tx, session)
		if txErr != nil {
			return txErr
		}
		task.Duration += duration
		return bknd.models.Tasks.UpdateTx(tx, task)
	})
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			bknd.errResourceNotFound(w, r)
		default:
			bknd.errInternalServerError(w, r, err)
		}
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/session/%s", session.ID))

	err = bknd.writeJSON(w, http.StatusCreated, wrapper{"session": session}, headers)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
	}
}

func (bknd *backend) showSessionsForTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bknd.readIdParam(r, "id")
	if err != nil {
		bknd.errBadRequest(w, r, err)
		return
	}
	sessions, err := bknd.models.Session.GetForTask(id)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
		return
	}
	err = bknd.writeJSON(w, http.StatusOK, wrapper{"sessions": sessions}, nil)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
	}
}
