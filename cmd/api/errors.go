package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	errorResponse(w, http.StatusInternalServerError, "the server encountered a problem and could not process your request")
}

func (app *application) forbiddenErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("forbidden error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	errorResponse(w, http.StatusForbidden, "you do not have permission to access this resource")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	errorResponse(w, http.StatusBadRequest, err.Error())
}

func (app *application) conflicResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("conflict error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	errorResponse(w, http.StatusConflict, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("not foundr", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	errorResponse(w, http.StatusNotFound, err.Error())
}

func (app *application) notContent(w http.ResponseWriter) {
	noContentResponse(w, http.StatusNoContent)
}

func (app *application) unauthorizeErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("unauthorize error response ", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	errorResponse(w, http.StatusUnauthorized, err.Error())
}

func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("unauthorize error response ", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	errorResponse(w, http.StatusUnauthorized, err.Error())
}
