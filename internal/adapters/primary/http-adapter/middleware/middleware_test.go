package middleware_test

import (
	"credit-service/internal/adapters/primary/http-adapter/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	called := false

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusCreated)
	})

	loggedHandler := middleware.LoggingMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	loggedHandler.ServeHTTP(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()

	if !called {
		t.Errorf("expected handler to be called")
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}
}
