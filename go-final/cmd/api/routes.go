package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/courses", app.listCoursesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/courses", app.requireActivatedUser(app.createCourseHandler))
	router.HandlerFunc(http.MethodGet, "/v1/courses/:id", app.showCourseHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/courses/:id", app.requireActivatedUser(app.updateCourseHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/courses/:id", app.requireActivatedUser(app.deleteCourseHandler))

	router.HandlerFunc(http.MethodPost, "/v1/instructors", app.requireActivatedUser(app.createInstructorHandler))
	router.HandlerFunc(http.MethodGet, "/v1/instructors", app.listInstructorsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/instructors/:id", app.showInstructorHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/instructors/:id", app.requireActivatedUser(app.updateInstructorHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/instructors/:id", app.requireActivatedUser(app.deleteInstructorHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.HandlerFunc(http.MethodPost, "/v1/files/send", app.fileHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
