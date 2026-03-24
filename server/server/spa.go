package server

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// webappFS holds the built webapp assets.
// The dist/ directory is copied into server/server/ by the Makefile before building.
//
//go:embed all:dist
var webappFS embed.FS

// serveSPA configures the SPA handler on the router.
// Static assets from dist/ are served directly. All other paths fall through
// to index.html for client-side routing.
func (s *SkyBookServer) serveSPA(router *mux.Router) {
	distFS, err := fs.Sub(webappFS, "dist")
	if err != nil {
		s.log.Error("Failed to access embedded webapp", "error", err)
		return
	}

	fileServer := http.FileServer(http.FS(distFS))

	spaHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the exact file first
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		// Check if the file exists in the embedded FS
		f, err := distFS.Open(strings.TrimPrefix(path, "/"))
		if err == nil {
			f.Close()

			// Hashed assets get long-lived cache headers;
			// index.html must never be cached (so new deploys are picked up).
			if strings.HasPrefix(path, "/assets/") {
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			} else if path == "/index.html" {
				w.Header().Set("Cache-Control", "no-cache")
			}

			fileServer.ServeHTTP(w, r)
			return
		}

		// File not found — serve index.html for client-side routing
		w.Header().Set("Cache-Control", "no-cache")
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})

	// Register as both a PathPrefix catch-all AND the NotFoundHandler
	// to ensure all unmatched routes fall through to the SPA.
	router.PathPrefix("/").Handler(spaHandler)
	router.NotFoundHandler = spaHandler
}
