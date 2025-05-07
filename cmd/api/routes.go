package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (bknd *backend) routes() http.Handler {
	router := httprouter.New()

	router.MethodNotAllowed = http.HandlerFunc(bknd.errMethodNotAllowed)
	router.NotFound = http.HandlerFunc(bknd.errResourceNotFound)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", bknd.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/tasks/:id", bknd.showTaskHandler)

	router.HandlerFunc(http.MethodPost, "/v1/categories", bknd.createCategoriesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/categories/:id", bknd.showCategoriesHandler)

	return bknd.recoverPanic(router)
}
