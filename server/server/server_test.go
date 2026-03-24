package server

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/root-gg/skybook/common"
	"github.com/root-gg/skybook/metadata"
)

func TestNewSkyBookServer(t *testing.T) {
	config := common.NewConfig()
	backend, err := metadata.NewBackend(config.Database, slog.Default())
	if err != nil {
		t.Fatalf("failed to init backend for test: %v", err)
	}
	defer backend.Shutdown()

	srv := NewSkyBookServer(config, backend, slog.Default())
	if srv == nil {
		t.Fatalf("expected server instance to be created")
	}
	if srv.router == nil {
		t.Errorf("expected router to be initialized")
	}
	if srv.server == nil {
		t.Errorf("expected underlying HTTP server to be initialized")
	}
}

func TestSetupRouter_APIRoutes(t *testing.T) {
	config := common.NewConfig()
	backend, err := metadata.NewBackend(config.Database, slog.Default())
	if err != nil {
		t.Fatalf("failed to init backend for test: %v", err)
	}
	defer backend.Shutdown()

	srv := NewSkyBookServer(config, backend, slog.Default())
	router := srv.router

	tests := []struct {
		method string
		path   string
		expect int
	}{
		{"GET", "/health", http.StatusOK},
		{"GET", "/api/v1/config", http.StatusOK},
		{"GET", "/api/v1/jumps", http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tc.expect {
				t.Errorf("expected status %d for %s, got %d", tc.expect, tc.path, rec.Code)
			}
		})
	}
}

func TestServerStartShutdown(t *testing.T) {
	config := common.NewConfig()
	// Use port 0 to let the OS pick an available port to avoid conflicts
	config.Server.ListenPort = 0

	backend, err := metadata.NewBackend(config.Database, slog.Default())
	if err != nil {
		t.Fatalf("failed to init backend for test: %v", err)
	}
	defer backend.Shutdown()

	srv := NewSkyBookServer(config, backend, slog.Default())

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	// Give it a moment to start up
	time.Sleep(100 * time.Millisecond)

	err = srv.Shutdown()
	if err != nil {
		t.Errorf("failed to gracefully shutdown server: %v", err)
	}

	// Verify Start() returns nil or http.ErrServerClosed smoothly
	serveErr := <-errCh
	if serveErr != nil {
		t.Errorf("expected pure shutdown without error, got: %v", serveErr)
	}
}
