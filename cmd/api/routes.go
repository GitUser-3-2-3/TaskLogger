package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

func (bknd *backend) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", bknd.healthcheckHandler)

	return router
}

func (bknd *backend) healthcheckHandler(w http.ResponseWriter, _ *http.Request) {
	log.Info().Msg("health check requested")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := `{"status":"available","environment":"` + bknd.config.env + `"}`
	_, _ = w.Write([]byte(response))
}
