package handlers

import (
	"net/http"

	"github.com/root-gg/skybook/common"
)

// HealthHandler returns a simple health check response.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	common.WriteJSON(w, map[string]string{"status": "ok"}, http.StatusOK)
}

// ConfigHandler returns a frontend-safe subset of the server configuration.
func ConfigHandler(config *common.Config) http.HandlerFunc {
	type publicConfig struct {
		UnitSystem      string `json:"unitSystem"`
		DefaultJumpType string `json:"defaultJumpType"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		common.WriteJSON(w, publicConfig{
			UnitSystem:      config.Defaults.UnitSystem,
			DefaultJumpType: config.Defaults.DefaultJumpType,
		}, http.StatusOK)
	}
}
