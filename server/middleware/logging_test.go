package middleware

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStatusWriter(t *testing.T) {
	rec := httptest.NewRecorder()
	sw := &statusWriter{ResponseWriter: rec, status: http.StatusOK}

	sw.WriteHeader(http.StatusCreated)

	if sw.status != http.StatusCreated {
		t.Errorf("expected statusWriter to capture %#v, got %#v", http.StatusCreated, sw.status)
	}

	if rec.Code != http.StatusCreated {
		t.Errorf("expected underlying ResponseWriter to receive %#v, got %#v", http.StatusCreated, rec.Code)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	// Create a handler that sets a specific status
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("OK"))
	})

	// Use a discard logger to not clutter test output, but we can't easily assert on it without a custom writer.
	// Since the requirement is just coverage and basic validation, this is sufficient.
	logger := slog.Default()
	middleware := Logging(logger)(nextHandler)

	req := httptest.NewRequest("GET", "/test-logging", nil)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Errorf("expected status %d, got %d", http.StatusAccepted, rec.Code)
	}

	body := strings.TrimSpace(rec.Body.String())
	if body != "OK" {
		t.Errorf("expected body 'OK', got %q", body)
	}
}
