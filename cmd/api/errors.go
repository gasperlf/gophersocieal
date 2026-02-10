package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	err = errorResponse(w, http.StatusInternalServerError, "the server encountered a problem and could not process your request")
	if err != nil {
		app.logger.Errorw("failed to send internal server error response", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	}
}

func (app *application) forbiddenErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("forbidden error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	err = errorResponse(w, http.StatusForbidden, "you do not have permission to access this resource")
	if err != nil {
		app.logger.Errorw("failed to send forbidden error response", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	}
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	err = errorResponse(w, http.StatusBadRequest, err.Error())
	if err != nil {
		app.logger.Errorw("failed to send bad request response", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	}
}

func (app *application) conflicResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("conflict error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	err = errorResponse(w, http.StatusConflict, err.Error())
	if err != nil {
		app.logger.Errorw("failed to send conflict error response", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	}
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("not foundr", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	err = errorResponse(w, http.StatusNotFound, err.Error())
	if err != nil {
		app.logger.Errorw("failed to send not found response", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	}
}

func (app *application) notContent(w http.ResponseWriter) {
	err := noContentResponse(w, http.StatusNoContent)
	if err != nil {
		app.logger.Errorw("failed to send no content response", "error", err.Error())
	}
}

func (app *application) unauthorizeErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("unauthorize error response ", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	err = errorResponse(w, http.StatusUnauthorized, err.Error())
	if err != nil {
		app.logger.Errorw("failed to send unauthorized error response", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	}
}

func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("unauthorize error response ", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	err = errorResponse(w, http.StatusUnauthorized, err.Error())
	if err != nil {
		app.logger.Errorw("failed to send unauthorized basic error response", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	}
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path)
	w.Header().Set("Retry-After", retryAfter) // In a real implementation, this should be dynamic based on the rate limiter's state
	writeJSON(w, http.StatusTooManyRequests, "rate limit exceeded, please retry after "+retryAfter)
}
