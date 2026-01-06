package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"ontopsolutions.net/gasperlf/social/internal/store"
)

type userKey string

const contextKeyUser userKey = "user"

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	user := getUserFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {

	var request FollowUser
	if err := readJSON(w, r, &request); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(request); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	followerUser := getUserFromContext(r)

	ctx := r.Context()
	err := app.store.Followers.Follow(ctx, followerUser.ID, request.UserID)

	if err != nil {

		switch {
		case errors.Is(err, store.ErrorConflict):
			app.conflicResponse(w, r, err)
			return
		default:
			app.badRequestResponse(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {

	var request FollowUser
	if err := readJSON(w, r, &request); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(request); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	unfollowerUser := getUserFromContext(r)
	ctx := r.Context()

	err := app.store.Followers.Unfollow(ctx, unfollowerUser.ID, request.UserID)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := getParamAsInt(r, "userID")
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetByID(ctx, userID)
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

		ctx = context.WithValue(ctx, contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(contextKeyUser).(*store.User)
	return user
}

func getParamAsInt(r *http.Request, param string) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, param), 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
