package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/root-gg/skybook/common"
)

// Recovery returns middleware that recovers from panics and returns a 500 JSON response.
func Recovery(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("Panic recovered",
						"error", rec,
						"stack", string(debug.Stack()),
						"method", r.Method,
						"path", r.URL.Path,
					)
					common.WriteError(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
