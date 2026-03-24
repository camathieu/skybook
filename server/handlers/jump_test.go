package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"github.com/root-gg/skybook/common"
	"github.com/root-gg/skybook/handlers"
	"github.com/root-gg/skybook/metadata"
)

// testDB creates a fresh in-memory-like SQLite database for handler tests.
func testDB(t *testing.T) *metadata.Backend {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	db, err := metadata.NewBackend(common.DatabaseConfig{Path: dbPath}, slog.Default())
	if err != nil {
		t.Fatalf("NewBackend: %v", err)
	}
	t.Cleanup(func() { db.Shutdown() })
	return db
}

func jsonBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return bytes.NewBuffer(b)
}

func mustDecode[T any](t *testing.T, body *bytes.Buffer) T {
	t.Helper()
	var v T
	if err := json.NewDecoder(body).Decode(&v); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return v
}

// --- ListJumps ---

func TestListJumps_Empty(t *testing.T) {
	db := testDB(t)
	req := httptest.NewRequest("GET", "/api/v1/jumps", nil)
	rr := httptest.NewRecorder()
	handlers.ListJumps(db)(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	type resp struct {
		Items []any `json:"items"`
		Total int   `json:"total"`
	}
	got := mustDecode[resp](t, rr.Body)
	if got.Total != 0 {
		t.Errorf("expected total 0, got %d", got.Total)
	}
	if got.Items == nil {
		t.Error("items must be an empty array, not null")
	}
}

func TestListJumps_Pagination(t *testing.T) {
	db := testDB(t)
	for i := 0; i < 10; i++ {
		db.CreateJump(&common.Jump{
			UserID:   1,
			Date:     time.Now(),
			Dropzone: fmt.Sprintf("DZ%d", i),
			JumpType: common.JumpTypeFF,
		})
	}

	req := httptest.NewRequest("GET", "/api/v1/jumps?page=2&per_page=3&sort=number&order=asc", nil)
	rr := httptest.NewRecorder()
	handlers.ListJumps(db)(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	type resp struct {
		Items      []*common.Jump `json:"items"`
		Total      int            `json:"total"`
		TotalPages int            `json:"totalPages"`
	}
	got := mustDecode[resp](t, rr.Body)
	if got.Total != 10 {
		t.Errorf("expected total 10, got %d", got.Total)
	}
	if len(got.Items) != 3 {
		t.Errorf("expected 3 items (page 2), got %d", len(got.Items))
	}
	if got.TotalPages != 4 {
		t.Errorf("expected 4 total pages, got %d", got.TotalPages)
	}
}

func TestListJumps_PerPageClamp(t *testing.T) {
	db := testDB(t)
	req := httptest.NewRequest("GET", "/api/v1/jumps?per_page=999", nil)
	rr := httptest.NewRecorder()
	handlers.ListJumps(db)(rr, req)
	// Should still return 200 (clamped silently)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 (per_page clamped), got %d", rr.Code)
	}
}

func TestListJumps_InvalidSort(t *testing.T) {
	db := testDB(t)
	req := httptest.NewRequest("GET", "/api/v1/jumps?sort=injected_sql", nil)
	rr := httptest.NewRecorder()
	handlers.ListJumps(db)(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid sort, got %d", rr.Code)
	}
}

// --- CreateJump ---

func TestCreateJump_Append(t *testing.T) {
	db := testDB(t)
	body := jsonBody(t, map[string]any{
		"date":     time.Now().Format(time.RFC3339),
		"dropzone": "Skydive DeLand",
		"jumpType": "FF",
	})
	req := httptest.NewRequest("POST", "/api/v1/jumps", body)
	rr := httptest.NewRecorder()
	handlers.CreateJump(db)(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}
	got := mustDecode[common.Jump](t, rr.Body)
	if got.Number != 1 {
		t.Errorf("expected number 1, got %d", got.Number)
	}
}

func TestCreateJump_MissingDate(t *testing.T) {
	db := testDB(t)
	body := jsonBody(t, map[string]any{
		"dropzone": "DZ",
		"jumpType": "FF",
	})
	req := httptest.NewRequest("POST", "/api/v1/jumps", body)
	rr := httptest.NewRecorder()
	handlers.CreateJump(db)(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing date, got %d", rr.Code)
	}
}

func TestCreateJump_InvalidJumpType(t *testing.T) {
	db := testDB(t)
	body := jsonBody(t, map[string]any{
		"date":     time.Now().Format(time.RFC3339),
		"dropzone": "DZ",
		"jumpType": "INVALID",
	})
	req := httptest.NewRequest("POST", "/api/v1/jumps", body)
	rr := httptest.NewRecorder()
	handlers.CreateJump(db)(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid jump type, got %d", rr.Code)
	}
}

func TestCreateJump_InsertAt(t *testing.T) {
	db := testDB(t)
	// Create 3 jumps
	for i := 0; i < 3; i++ {
		db.CreateJump(&common.Jump{UserID: 1, Date: time.Now(), Dropzone: "DZ", JumpType: common.JumpTypeFF})
	}
	// Insert at position 2
	body := jsonBody(t, map[string]any{
		"date":     time.Now().Format(time.RFC3339),
		"dropzone": "Inserted",
		"jumpType": "WINGSUIT",
		"number":   2,
	})
	req := httptest.NewRequest("POST", "/api/v1/jumps", body)
	rr := httptest.NewRecorder()
	handlers.CreateJump(db)(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}
	got := mustDecode[common.Jump](t, rr.Body)
	if got.Number != 2 {
		t.Errorf("expected number 2, got %d", got.Number)
	}
}

// --- GetJump ---

func getJumpRequest(t *testing.T, db *metadata.Backend, id uint) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/jumps/%d", id), nil)
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", id)})
	rr := httptest.NewRecorder()
	handlers.GetJump(db)(rr, req)
	return rr
}

func TestGetJump_Found(t *testing.T) {
	db := testDB(t)
	j := &common.Jump{UserID: 1, Date: time.Now(), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	db.CreateJump(j)

	rr := getJumpRequest(t, db, j.ID)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestGetJump_NotFound(t *testing.T) {
	db := testDB(t)
	rr := getJumpRequest(t, db, 9999)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

// --- UpdateJump ---

func TestUpdateJump_Fields(t *testing.T) {
	db := testDB(t)
	j := &common.Jump{UserID: 1, Date: time.Now(), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	db.CreateJump(j)

	body := jsonBody(t, map[string]any{
		"date":     time.Now().Format(time.RFC3339),
		"dropzone": "Updated DZ",
		"jumpType": "WINGSUIT",
	})
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/jumps/%d", j.ID), body)
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", j.ID)})
	rr := httptest.NewRecorder()
	handlers.UpdateJump(db)(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	got := mustDecode[common.Jump](t, rr.Body)
	if got.Dropzone != "Updated DZ" {
		t.Errorf("expected dropzone 'Updated DZ', got %q", got.Dropzone)
	}
}

func TestUpdateJump_NotFound(t *testing.T) {
	db := testDB(t)
	body := jsonBody(t, map[string]any{
		"date":     time.Now().Format(time.RFC3339),
		"dropzone": "DZ",
		"jumpType": "FF",
	})
	req := httptest.NewRequest("PUT", "/api/v1/jumps/9999", body)
	req = mux.SetURLVars(req, map[string]string{"id": "9999"})
	rr := httptest.NewRecorder()
	handlers.UpdateJump(db)(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestUpdateJump_Move(t *testing.T) {
	db := testDB(t)
	for i := 0; i < 3; i++ {
		db.CreateJump(&common.Jump{UserID: 1, Date: time.Now(), Dropzone: "DZ", JumpType: common.JumpTypeFF})
	}
	// Move jump #3 to position #1
	j, _ := db.GetJumpByNumber(1, 3)
	body := jsonBody(t, map[string]any{
		"date":     j.Date.Format(time.RFC3339),
		"dropzone": j.Dropzone,
		"jumpType": string(j.JumpType),
		"number":   1,
	})
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/jumps/%d", j.ID), body)
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", j.ID)})
	rr := httptest.NewRecorder()
	handlers.UpdateJump(db)(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	got := mustDecode[common.Jump](t, rr.Body)
	if got.Number != 1 {
		t.Errorf("expected number 1 after move, got %d", got.Number)
	}
}

// --- DeleteJump ---

func TestDeleteJump_Success(t *testing.T) {
	db := testDB(t)
	j := &common.Jump{UserID: 1, Date: time.Now(), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	db.CreateJump(j)

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/jumps/%d", j.ID), nil)
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", j.ID)})
	rr := httptest.NewRecorder()
	handlers.DeleteJump(db)(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
	// Confirm it's gone
	rr2 := getJumpRequest(t, db, j.ID)
	if rr2.Code != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", rr2.Code)
	}
}

func TestDeleteJump_NotFound(t *testing.T) {
	db := testDB(t)
	req := httptest.NewRequest("DELETE", "/api/v1/jumps/9999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "9999"})
	rr := httptest.NewRecorder()
	handlers.DeleteJump(db)(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

// --- Autocomplete ---

func TestAutocomplete_Dropzone(t *testing.T) {
	db := testDB(t)
	for i := 0; i < 3; i++ {
		db.CreateJump(&common.Jump{
			UserID:   1,
			Date:     time.Now(),
			Dropzone: "Skydive DeLand",
			JumpType: common.JumpTypeFF,
		})
	}
	db.CreateJump(&common.Jump{
		UserID:   1,
		Date:     time.Now(),
		Dropzone: "Perris",
		JumpType: common.JumpTypeFF,
	})

	req := httptest.NewRequest("GET", "/api/v1/jumps/autocomplete/dropzone", nil)
	req = mux.SetURLVars(req, map[string]string{"field": "dropzone"})
	rr := httptest.NewRecorder()
	handlers.Autocomplete(db)(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var results []metadata.AutocompleteResult
	if err := json.NewDecoder(rr.Body).Decode(&results); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected autocomplete results")
	}
	// Top result should be "Skydive DeLand" with count 3
	if results[0].Value != "Skydive DeLand" {
		t.Errorf("expected top result 'Skydive DeLand', got %q", results[0].Value)
	}
}

func TestAutocomplete_InvalidField(t *testing.T) {
	db := testDB(t)
	req := httptest.NewRequest("GET", "/api/v1/jumps/autocomplete/hacked", nil)
	req = mux.SetURLVars(req, map[string]string{"field": "hacked"})
	rr := httptest.NewRecorder()
	handlers.Autocomplete(db)(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}
