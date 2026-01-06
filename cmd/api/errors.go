package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	errorResponse(w, http.StatusInternalServerError, "the server encountered a problem and could not process your request")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request: %s path: %s error: %s", r.Method, r.URL.Path, err)
	errorResponse(w, http.StatusBadRequest, err.Error())
}

func (app *application) conflicResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("conflict error: %s path: %s error: %s", r.Method, r.URL.Path, err)
	errorResponse(w, http.StatusConflict, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found: %s path: %s error: %s", r.Method, r.URL.Path, err)
	errorResponse(w, http.StatusNotFound, err.Error())
}

func (app *application) notContent(w http.ResponseWriter, r *http.Request) {
	log.Printf("not content: %s path: %s", r.Method, r.URL.Path)
	noContentResponse(w, http.StatusNoContent)
}
