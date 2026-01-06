package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"ontopsolutions.net/gasperlf/social/internal/store"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getParamAsInt(r, "userID")
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	user, err := app.store.Users.GetByID(r.Context(), userID)
	if err != nil {
		switch err {
		case store.ErrorNotFound:
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	err = app.jsonResponse(w, http.StatusOK, user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func getParamAsInt(r *http.Request, param string) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, param), 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
