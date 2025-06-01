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

	bknd.registerSessionRoutes(router)

	return bknd.recoverPanic(bknd.rateLimiter(router))
}

func (bknd *backend) registerCategoryRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/v1/categories", bknd.createCategoriesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/categories/:categoryId", bknd.showCategoriesHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/categories/:categoryId", bknd.updateCategoryHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/categories/:categoryId", bknd.deleteCategoryHandler)
	router.HandlerFunc(http.MethodGet, "/v1/categories/:categoryId/tasks", bknd.showTasksByCategory)
}

func (bknd *backend) registerUserRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet,
		"/v1/users/:userId/categories", bknd.showCategoriesForUserIdHandler)
	router.HandlerFunc(http.MethodGet,
		"/v1/users/:userId/categories/tasks", bknd.showCategoryDetailsForUserIdHandler)
}

func (bknd *backend) registerTaskRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/v1/tasks", bknd.createTaskHandler)
	router.HandlerFunc(http.MethodGet, "/v1/tasks/:taskId", bknd.showTaskHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/tasks/:taskId", bknd.updateTaskHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/tasks/:taskId", bknd.deleteTaskHandler)

	router.HandlerFunc(http.MethodGet, "/v1/tasks/:taskId/sessions", bknd.showSessionsForTaskHandler)
}

func (bknd *backend) registerSessionRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/v1/sessions", bknd.createSessionHandler)
}
