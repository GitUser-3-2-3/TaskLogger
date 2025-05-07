package main

import (
	"TaskLogger/internal/data"
	"net/http"
	"time"
)

func (bknd *backend) showTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bknd.readIdParam(r)
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
