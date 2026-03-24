package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/root-gg/skybook/common"
	"github.com/root-gg/skybook/handlers"
	"github.com/root-gg/skybook/metadata"
	"github.com/root-gg/skybook/middleware"
)

// SkyBookServer is the main HTTP server.
type SkyBookServer struct {
	config  *common.Config
	backend *metadata.Backend
	router  *mux.Router
	server  *http.Server
	log     *slog.Logger
}

// NewSkyBookServer creates a new server instance with routes and middleware.
func NewSkyBookServer(config *common.Config, backend *metadata.Backend, log *slog.Logger) *SkyBookServer {
	s := &SkyBookServer{
		config:  config,
		backend: backend,
		log:     log,
	}

	s.router = s.setupRouter()
	s.server = &http.Server{
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
}

// setupRouter creates the mux router with all routes and middleware.
func (s *SkyBookServer) setupRouter() *mux.Router {
	r := mux.NewRouter()

	// Global middleware chain
	r.Use(middleware.Recovery(s.log))
	r.Use(middleware.RequestID())
	r.Use(middleware.Logging(s.log))

	// Health check
	r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")

	// API v1
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/config", handlers.ConfigHandler(s.config)).Methods("GET")

	// Jumps — autocomplete must be registered before /{id} to avoid mux conflicts
	api.HandleFunc("/jumps/autocomplete/{field}", handlers.Autocomplete(s.backend)).Methods("GET")
	api.HandleFunc("/jumps", handlers.ListJumps(s.backend)).Methods("GET")
	api.HandleFunc("/jumps", handlers.CreateJump(s.backend)).Methods("POST")
	api.HandleFunc("/jumps/{id:[0-9]+}", handlers.GetJump(s.backend)).Methods("GET")
	api.HandleFunc("/jumps/{id:[0-9]+}", handlers.UpdateJump(s.backend)).Methods("PUT")
	api.HandleFunc("/jumps/{id:[0-9]+}", handlers.DeleteJump(s.backend)).Methods("DELETE")

	return r
}

// Start begins listening for HTTP requests.
func (s *SkyBookServer) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.ListenAddress, s.config.Server.ListenPort)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}

	s.log.Info("Server listening", "address", listener.Addr().String())

	if err := s.server.Serve(listener); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("serve: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the HTTP server.
func (s *SkyBookServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.log.Info("Shutting down HTTP server")
	return s.server.Shutdown(ctx)
}
