package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (bknd *backend) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", bknd.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/tasks/:id", bknd.showTaskHandler)

	return router
}
