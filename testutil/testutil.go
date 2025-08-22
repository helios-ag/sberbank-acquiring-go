package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// NewTestServer creates a new test server with the given handler function
func NewTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, func()) {
	s := httptest.NewServer(handler)
	t.Cleanup(s.Close)
	return s, s.Close
}

// WriteJSON writes a JSON response with the given status code
func WriteJSON(w http.ResponseWriter, statusCode int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(v)
}

// MustJSON marshals the given value to JSON, panicking on error
func MustJSON(t *testing.T, v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	return data
}

// ErrorResponse represents a standard error response from the API
type ErrorResponse struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}
