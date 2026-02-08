package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"ontopsolutions.net/gasperlf/social/internal/store"
)

func (app *application) BasicMidleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("unauthorization header is missing"))
				return
			}
			//parse it -> get the base 64
			parts := strings.Split(authHeader, "")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("unauthorization header is malformed"))
				return
			}
			// decode it
			decode, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicErrorResponse(w, r, err)
				return
			}

			//check the credentials
			username := app.config.auth.basic.user
			pass := app.config.auth.basic.pass

			creds := strings.SplitN(string(decode), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//read the auth header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("unauthorization header is missing"))
			return
		}

		//parse it -> get the base 64
		parts := strings.Split(authHeader, "")

		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizeErrorResponse(w, r, fmt.Errorf("unauthorization header is malformed"))
			return
		}

		token := parts[1]
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizeErrorResponse(w, r, err)
			return
		}
		//add the claims to the context
		claims, _ := jwtToken.Claims.(jwt.MapClaims)
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64) //claims["sub"].(string)
		if err != nil {
			app.unauthorizeErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			app.unauthorizeErrorResponse(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func (app *application) CheckPostOwnership(roleRequired string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// if it is the owner of the post
		user := getUserFromContext(r)
		post := getPostFromContext(r)

		if user.ID == post.UserID {
			next.ServeHTTP(w, r)
			return
		}

		//role precedence check
		allowance, err := app.checkRolePrecedence(r.Context(), user, roleRequired)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		if !allowance {
			app.forbiddenErrorResponse(w, r, fmt.Errorf("you don't have the required permissions to perform this action"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {

	// find the role of the user
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	// find the precedence of the role required
	return user.Role.Level >= role.Level, nil
}
