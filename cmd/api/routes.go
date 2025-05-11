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

	bknd.registerCategoryRoutes(router)

	bknd.registerUserRoutes(router)

	bknd.registerTaskRoutes(router)

	return bknd.recoverPanic(bknd.rateLimiter(router))
}

func (bknd *backend) registerCategoryRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/v1/categories", bknd.createCategoriesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/categories/:id", bknd.showCategoriesHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/categories/:id", bknd.updateCategoryHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/categories/:id", bknd.deleteCategoryHandler)
	router.HandlerFunc(http.MethodGet, "/v1/categories/:id/tasks", bknd.showTasksByCategory)
}

func (bknd *backend) registerUserRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet,
		"/v1/users/:userId/categories", bknd.showCategoriesForUserIdHandler)
	router.HandlerFunc(http.MethodGet,
		"/v1/users/:userId/categories/tasks", bknd.showCategoriesDetailsForUserIdHandler)
}

func (bknd *backend) registerTaskRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/v1/tasks", bknd.createTaskHandler)
	router.HandlerFunc(http.MethodGet, "/v1/tasks/:id", bknd.showTaskHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/tasks/:id", bknd.updateTaskHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/tasks/:id", bknd.deleteTaskHandler)
}
