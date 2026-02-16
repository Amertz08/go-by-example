package rest

import (
	"io"
	"net/http"
	"testing"
)

// newRequest creates a new HTTP request for testing
func newRequest(t testing.TB, method, target string, body io.Reader) *http.Request {
	t.Helper()
	req, err := http.NewRequest(method, target, body)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	return req
}
