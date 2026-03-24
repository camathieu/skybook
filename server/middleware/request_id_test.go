package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestIDMiddleware(t *testing.T) {
	var capturedID string

	// Create a handler that records the injected RequestID from the context
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = GetRequestID(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	middleware := RequestID()(nextHandler)

	req := httptest.NewRequest("GET", "/test-request-id", nil)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	// Verify the context value was set and captured
	if capturedID == "" {
		t.Errorf("expected non-empty RequestID injected into context")
	}

	// Verify the response header is exactly the same ID
	headerID := rec.Header().Get("X-Request-ID")
	if headerID != capturedID {
		t.Errorf("expected 'X-Request-ID' header to be %q, got %q", capturedID, headerID)
	}
}

func TestGetRequestID_Empty(t *testing.T) {
	// Context without the requestIDKey
	ctx := context.Background()

	id := GetRequestID(ctx)
	if id != "" {
		t.Errorf("expected empty string when no request ID is in context, got %q", id)
	}
}
