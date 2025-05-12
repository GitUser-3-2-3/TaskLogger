package main

import (
	"TaskLogger/internal/data"
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
		TaskID       int64            `json:"task_id"`
		SessionStart time.Time        `json:"session_start"`
		SessionEnd   time.Time        `json:"session_end"`
		Note         string           `json:"note"`
		SessionType  data.SessionType `json:"session_type"`
	}
	err := bknd.readJSON(w, r, &input)
	if err != nil {
		bknd.errBadRequest(w, r, err)
		return
	}
	if input.SessionEnd.Before(input.SessionStart) {
		bknd.errBadRequest(w, r, fmt.Errorf("session start must be before session end"))
		return
	}
	duration := int(input.SessionEnd.Sub(input.SessionStart).Minutes())
	session := &data.Session{
		TaskID:       input.TaskID,
		SessionStart: input.SessionStart,
		SessionEnd:   input.SessionEnd,
		Duration:     duration,
		Note:         input.Note,
		SessionType:  input.SessionType,
	}
	var sessionId int64

	err = bknd.withTransaction(func(tx *sql.Tx) error {
		var txErr error
		sessionId, txErr = bknd.models.Session.InsertTx(tx, session)
		if txErr != nil {
			return txErr
		}
		return bknd.updateTaskDuration(tx, session.TaskID, duration)
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
	headers.Set("Location", fmt.Sprintf("/v1/session/%d", sessionId))

	session.ID = sessionId
	err = bknd.writeJSON(w, http.StatusCreated, wrapper{"session": session}, headers)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
	}
}

func (bknd *backend) updateTaskDuration(tx *sql.Tx, taskId int64, duration int) error {
	task, err := bknd.models.Tasks.GetByIdTx(tx, taskId)
	if err != nil {
		return err
	}
	task.TotalDuration += duration
	return bknd.models.Tasks.UpdateTx(tx, task)
}
