package main

import (
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {

	app := newTestApplication(t)
	mux := mount(app)
	testToken, _ := app.authenticator.GenerateToken(nil)

	t.Run("should not allow unauthorized request", func(t *testing.T) {

		req, err := http.NewRequest("GET", "/v1/users/1", nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should allow authenticated requests", func(t *testing.T) {

		req, err := http.NewRequest("GET", "/v1/users/1", nil)

		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusOK, rr.Code)

	})
}
