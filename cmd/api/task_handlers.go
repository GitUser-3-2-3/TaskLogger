package main

import (
	"TaskLogger/internal/data"
	"net/http"
	"time"
)

func (bknd *backend) showTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bknd.readIdParam(r)
	if err != nil {
		http.NotFound(w, r)
		bknd.logger.Err(err).Msg("Error while reading id")
		return
	}
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
		Deadline:      time.Now().Add(24 * time.Hour),
		UserID:        1,
		CategoryID:    1,
	}
	err = bknd.writeJSON(w, http.StatusOK, task, nil)
	if err != nil {
		bknd.logger.Err(err).Msg("Error while writing to JSON")
		http.Error(w, "Couldn't process your request", http.StatusInternalServerError)
	}
}
