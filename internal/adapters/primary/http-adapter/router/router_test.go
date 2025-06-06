package router_test

import (
	"credit-service/internal/adapters/primary/http-adapter/router"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockHandler struct{}

func (m *mockHandler) Execute(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("execute called"))
}

func (m *mockHandler) Cache(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("cache called"))
}

func TestRouter_RegisterRoutes(t *testing.T) {
	r := router.NewRouter()
	h := &mockHandler{}
	r.RegisterRoutes(h)

	ts := httptest.NewServer(r.Router())
	defer ts.Close()

	tests := []struct {
		method       string
		path         string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			method:       http.MethodPost,
			path:         "/execute",
			body:         `{"key":"value"}`,
			expectedCode: http.StatusOK,
			expectedBody: "execute called",
		},
		{
			method:       http.MethodGet,
			path:         "/cache",
			expectedCode: http.StatusOK,
			expectedBody: "cache called",
		},
		{
			method:       http.MethodGet,
			path:         "/not-found",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, ts.URL+tt.path, strings.NewReader(tt.body))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("could not send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != tt.expectedCode {
			t.Errorf("expected status %d, got %d", tt.expectedCode, resp.StatusCode)
		}

		if tt.expectedBody != "" {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("could not read response body: %v", err)
			}
			bodyStr := string(bodyBytes)
			if bodyStr != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, bodyStr)
			}

		}
	}
}
