package main

import (
	"net/http"
)

func (bknd *backend) healthcheckHandler(w http.ResponseWriter, _ *http.Request) {
	bknd.logger.Info().Msg("health check requested")

	data := map[string]string{"status": "available",
		"version": version,
		"env":     "development",
	}
	err := bknd.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		bknd.logger.Error().Err(err).Msg("Error marshalling response")
		http.Error(w, "cannot process your request", http.StatusInternalServerError)
		return
	}
}
