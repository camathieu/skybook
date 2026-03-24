package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/root-gg/skybook/common"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()

	HealthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}

	expectedBody := `{"status":"ok"}`
	body := strings.TrimSpace(rec.Body.String())
	if body != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, body)
	}
}

func TestConfigHandler(t *testing.T) {
	cfg := common.NewConfig()
	cfg.Defaults.UnitSystem = "imperial"
	cfg.Defaults.DefaultJumpType = "FF"

	handler := ConfigHandler(cfg)

	req := httptest.NewRequest("GET", "/api/v1/config", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rec.Code)
	}

	var resp map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify the frontend-safe exposure behavior
	if resp["unitSystem"] != "imperial" {
		t.Errorf("expected unitSystem 'imperial', got %q", resp["unitSystem"])
	}
	if resp["defaultJumpType"] != "FF" {
		t.Errorf("expected defaultJumpType 'FF', got %q", resp["defaultJumpType"])
	}
}
