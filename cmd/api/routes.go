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

	router.HandlerFunc(http.MethodPost, "/v1/categories", bknd.createCategoriesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/categories/:id", bknd.showCategoriesHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/categories/:id", bknd.updateCategoryHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/:userId/categories", bknd.showCategoriesByUserIdHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/categories/:id", bknd.deleteCategoryHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tasks", bknd.createTaskHandler)
	router.HandlerFunc(http.MethodGet, "/v1/tasks/:id", bknd.showTaskHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/tasks/:id", bknd.deleteTaskHandler)

	return bknd.recoverPanic(bknd.rateLimiter(router))
}
