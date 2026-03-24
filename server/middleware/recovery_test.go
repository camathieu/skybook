package middleware

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecoveryMiddleware(t *testing.T) {
	// Create a handler that panics
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("intentional test panic")
	})

	// Wrap it with Recovery middleware
	logger := slog.Default()
	middleware := Recovery(logger)(panicHandler)

	req := httptest.NewRequest("GET", "/test-panic", nil)
	rec := httptest.NewRecorder()

	// This should NOT crash the test. It should recover and return 500.
	middleware.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500 on panic, got %d", rec.Code)
	}

	// Verify the JSON error response is structurally correct based on common.WriteError
	body := strings.TrimSpace(rec.Body.String())
	if !strings.Contains(body, "Internal Server Error") {
		t.Errorf("expected body to contain 'Internal Server Error', got %q", body)
	}
}

func TestRecoveryMiddleware_NoPanic(t *testing.T) {
	// Create a healthy handler
	healthyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("healthy"))
	})

	logger := slog.Default()
	middleware := Recovery(logger)(healthyHandler)

	req := httptest.NewRequest("GET", "/test-healthy", nil)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := strings.TrimSpace(rec.Body.String())
	if body != "healthy" {
		t.Errorf("expected body 'healthy', got %q", body)
	}
}
