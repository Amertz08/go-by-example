package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealth(t *testing.T) {
	req := newRequest(t, http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	expectedBody := "OK"
	if got := w.Body.String(); !strings.Contains(got, expectedBody) {
		t.Errorf("expected body %q, got %q", expectedBody, w.Body.String())
	}
}
