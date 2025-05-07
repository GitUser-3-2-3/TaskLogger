package main

import (
	"net/http"
)

func (bknd *backend) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	bknd.logger.Info().Msg("health check requested")

	data := map[string]string{"status": "available",
		"version": version,
		"env":     "development",
	}
	err := bknd.writeJSON(w, http.StatusOK, wrapper{"system_info": data}, nil)
	if err != nil {
		bknd.errInternalServerError(w, r, err)
		return
	}
}
