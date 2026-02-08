package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"ontopsolutions.net/gasperlf/social/internal/auth"
	"ontopsolutions.net/gasperlf/social/internal/store"
	"ontopsolutions.net/gasperlf/social/internal/store/cache"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()

	logger := zap.Must(zap.NewProduction()).Sugar()
	mockStore := store.NewMockStore()
	mockCacheUser := cache.NewMockCache()
	testAuth := &auth.TestAuthenticator{}

	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStore:    mockCacheUser,
		authenticator: testAuth,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	t.Helper()

	if expected != actual {
		t.Errorf("expected status %d, got %d", expected, actual)
	}
}
