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
	today := common.Today()
	for i := 0; i < 10; i++ {
		db.CreateJump(&common.Jump{
			UserID:   1,
			Date:     today,
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

func TestListJumps_DateFromFilter(t *testing.T) {
	db := testDB(t)
	// Create 3 jumps on different dates
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 1), Dropzone: "DZ", JumpType: common.JumpTypeFF})
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 5), Dropzone: "DZ", JumpType: common.JumpTypeFF})
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 10), Dropzone: "DZ", JumpType: common.JumpTypeFF})

	// Filter: only jumps on or after March 5
	req := httptest.NewRequest("GET", "/api/v1/jumps?date_from=2025-03-05&sort=number&order=asc", nil)
	rr := httptest.NewRecorder()
	handlers.ListJumps(db)(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	type resp struct {
		Items []*common.Jump `json:"items"`
		Total int64          `json:"total"`
	}
	got := mustDecode[resp](t, rr.Body)
	if got.Total != 2 {
		t.Errorf("expected 2 jumps (Mar 5 and Mar 10), got %d", got.Total)
	}
	for _, j := range got.Items {
		if j.Date.Time.Before(common.NewDateOnly(2025, time.March, 5).Time) {
			t.Errorf("jump %d has date %s which is before date_from 2025-03-05", j.Number, j.Date.DayString())
		}
	}
}

func TestListJumps_DateToFilter(t *testing.T) {
	db := testDB(t)
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 1), Dropzone: "DZ", JumpType: common.JumpTypeFF})
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 5), Dropzone: "DZ", JumpType: common.JumpTypeFF})
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 10), Dropzone: "DZ", JumpType: common.JumpTypeFF})

	// Filter: only jumps on or before March 5
	req := httptest.NewRequest("GET", "/api/v1/jumps?date_to=2025-03-05", nil)
	rr := httptest.NewRecorder()
	handlers.ListJumps(db)(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	type resp struct {
		Total int64 `json:"total"`
	}
	got := mustDecode[resp](t, rr.Body)
	if got.Total != 2 {
		t.Errorf("expected 2 jumps (Mar 1 and Mar 5), got %d", got.Total)
	}
}

func TestListJumps_InvalidDateFrom(t *testing.T) {
	db := testDB(t)
	req := httptest.NewRequest("GET", "/api/v1/jumps?date_from=not-a-date", nil)
	rr := httptest.NewRecorder()
	handlers.ListJumps(db)(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid date_from, got %d", rr.Code)
	}
}

// --- CreateJump ---

func TestCreateJump_Append(t *testing.T) {
	db := testDB(t)
	body := jsonBody(t, map[string]any{
		"date":     common.Today().DayString(),
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
		"date":     common.Today().DayString(),
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
	today := common.Today()
	// Create 3 jumps
	for i := 0; i < 3; i++ {
		db.CreateJump(&common.Jump{UserID: 1, Date: today, Dropzone: "DZ", JumpType: common.JumpTypeFF})
	}
	// Insert at position 2
	body := jsonBody(t, map[string]any{
		"date":     today.DayString(),
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
	j := &common.Jump{UserID: 1, Date: common.Today(), Dropzone: "DZ", JumpType: common.JumpTypeFF}
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
	j := &common.Jump{UserID: 1, Date: common.Today(), Dropzone: "DZ", JumpType: common.JumpTypeFF}
	db.CreateJump(j)

	body := jsonBody(t, map[string]any{
		"date":     common.Today().DayString(),
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
		"date":     common.Today().DayString(),
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
	today := common.Today()
	for i := 0; i < 3; i++ {
		db.CreateJump(&common.Jump{UserID: 1, Date: today, Dropzone: "DZ", JumpType: common.JumpTypeFF})
	}
	// Move jump #3 to position #1
	j, _ := db.GetJumpByNumber(1, 3)
	body := jsonBody(t, map[string]any{
		"date":     j.Date.DayString(),
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
	j := &common.Jump{UserID: 1, Date: common.Today(), Dropzone: "DZ", JumpType: common.JumpTypeFF}
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
	// Create 3 jumps at Skydive DeLand, then 1 more recent jump at Perris
	past := common.NewDateOnly(2025, time.March, 1)
	for i := 0; i < 3; i++ {
		db.CreateJump(&common.Jump{
			UserID:   1,
			Date:     past,
			Dropzone: "Skydive DeLand",
			JumpType: common.JumpTypeFF,
		})
	}
	// Perris was used more recently — should sort first
	db.CreateJump(&common.Jump{
		UserID:   1,
		Date:     common.NewDateOnly(2025, time.March, 8),
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
	var results []string
	if err := json.NewDecoder(rr.Body).Decode(&results); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected autocomplete results")
	}
	// Perris was used most recently — should sort first
	if results[0] != "Perris" {
		t.Errorf("expected top result 'Perris' (most recent), got %q", results[0])
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

func TestAutocomplete_AlphaSort(t *testing.T) {
	db := testDB(t)
	// Create jumps: Perris most recently but Empuriabrava alphabetically first
	// Dates must be in chronological order (date validation enforced)
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 1), Dropzone: "Empuriabrava", JumpType: common.JumpTypeFF})
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 5), Dropzone: "Perris", JumpType: common.JumpTypeFF})
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 8), Dropzone: "Skydive DeLand", JumpType: common.JumpTypeFF})

	req := httptest.NewRequest("GET", "/api/v1/jumps/autocomplete/dropzone?sort=alpha", nil)
	req = mux.SetURLVars(req, map[string]string{"field": "dropzone"})
	rr := httptest.NewRecorder()
	handlers.Autocomplete(db)(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var results []string
	if err := json.NewDecoder(rr.Body).Decode(&results); err != nil {
		t.Fatalf("decode: %v", err)
	}
	// Alphabetically: Empuriabrava < Perris < Skydive DeLand
	if len(results) < 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0] != "Empuriabrava" {
		t.Errorf("expected top result 'Empuriabrava' (alpha), got %q", results[0])
	}
	if results[1] != "Perris" {
		t.Errorf("expected second result 'Perris' (alpha), got %q", results[1])
	}
}

// --- Date Validation ---

func TestCreateJump_DateBeforePrevious(t *testing.T) {
	db := testDB(t)
	// Create a jump on March 8
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 8), Dropzone: "DZ", JumpType: common.JumpTypeFF})
	// Attempt to append a jump on March 5 (before last) — should fail
	body := jsonBody(t, map[string]any{
		"date":     "2025-03-05",
		"dropzone": "DZ",
		"jumpType": "FF",
	})
	req := httptest.NewRequest("POST", "/api/v1/jumps", body)
	rr := httptest.NewRecorder()
	handlers.CreateJump(db)(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for date before previous, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestCreateJump_SameDayValid(t *testing.T) {
	db := testDB(t)
	today := common.NewDateOnly(2025, time.March, 8)
	// Create 3 jumps on the same day — all should be valid
	for i := 0; i < 3; i++ {
		body := jsonBody(t, map[string]any{
			"date":     today.DayString(),
			"dropzone": "DZ",
			"jumpType": "FF",
		})
		req := httptest.NewRequest("POST", "/api/v1/jumps", body)
		rr := httptest.NewRecorder()
		handlers.CreateJump(db)(rr, req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("jump %d: expected 201, got %d: %s", i+1, rr.Code, rr.Body.String())
		}
	}
}

func TestUpdateJump_DateBreaksOrder(t *testing.T) {
	db := testDB(t)
	// Create 3 jumps: March 5, March 8, March 10
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 5), Dropzone: "DZ", JumpType: common.JumpTypeFF})
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 8), Dropzone: "DZ", JumpType: common.JumpTypeFF})
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 10), Dropzone: "DZ", JumpType: common.JumpTypeFF})

	// Try to update jump #2 to March 3 (before #1's March 5) — should fail
	j, _ := db.GetJumpByNumber(1, 2)
	body := jsonBody(t, map[string]any{
		"date":     "2025-03-03",
		"dropzone": "DZ",
		"jumpType": "FF",
	})
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/jumps/%d", j.ID), body)
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", j.ID)})
	rr := httptest.NewRecorder()
	handlers.UpdateJump(db)(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for date out of order, got %d: %s", rr.Code, rr.Body.String())
	}
}

// --- API Date Format ---

func TestCreateJump_ReturnsDateOnly(t *testing.T) {
	db := testDB(t)
	body := jsonBody(t, map[string]any{
		"date":     "2025-03-08",
		"dropzone": "DZ",
		"jumpType": "FF",
	})
	req := httptest.NewRequest("POST", "/api/v1/jumps", body)
	rr := httptest.NewRecorder()
	handlers.CreateJump(db)(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}

	// Verify the JSON response has date in YYYY-MM-DD format
	var raw map[string]any
	json.NewDecoder(rr.Body).Decode(&raw)
	dateStr, ok := raw["date"].(string)
	if !ok {
		t.Fatal("expected date field in response")
	}
	if dateStr != "2025-03-08" {
		t.Errorf("expected date '2025-03-08', got %q", dateStr)
	}
}

func TestCreateJump_AcceptsRFC3339(t *testing.T) {
	db := testDB(t)
	body := jsonBody(t, map[string]any{
		"date":     "2025-03-08T14:30:00Z",
		"dropzone": "DZ",
		"jumpType": "FF",
	})
	req := httptest.NewRequest("POST", "/api/v1/jumps", body)
	rr := httptest.NewRecorder()
	handlers.CreateJump(db)(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}

	// Verify the time component was stripped in the response
	var raw map[string]any
	json.NewDecoder(rr.Body).Decode(&raw)
	dateStr := raw["date"].(string)
	if dateStr != "2025-03-08" {
		t.Errorf("expected date '2025-03-08' (time stripped), got %q", dateStr)
	}
}

func TestUpdateJump_MoveDateAndNumber(t *testing.T) {
	db := testDB(t)
	// Create 3 jumps: #1 (Mar 1), #2 (Mar 5), #3 (Mar 10)
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 1), Dropzone: "DZ", JumpType: common.JumpTypeFF})
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 5), Dropzone: "DZ", JumpType: common.JumpTypeFF})
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 10), Dropzone: "DZ", JumpType: common.JumpTypeFF})

	// Update jump #1 (Mar 1) to date=Mar 7, number=2
	// At current position (#1): Mar 7 > next (#2 Mar 5) → would be REJECTED by old code
	// At target position (#2): prev=#1 (which is itself, skip), next=#3 (Mar 10) → Mar 7 < Mar 10 → VALID
	j, _ := db.GetJumpByNumber(1, 1)
	body := jsonBody(t, map[string]any{
		"date":     "2025-03-07",
		"dropzone": "DZ",
		"jumpType": "FF",
		"number":   2,
	})
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/jumps/%d", j.ID), body)
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", j.ID)})
	rr := httptest.NewRecorder()
	handlers.UpdateJump(db)(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var result common.Jump
	json.NewDecoder(rr.Body).Decode(&result)
	if result.Number != 2 {
		t.Errorf("expected number 2, got %d", result.Number)
	}
	if result.Date.DayString() != "2025-03-07" {
		t.Errorf("expected date 2025-03-07, got %s", result.Date.DayString())
	}
}

// TestUpdateJump_MoveFailsDateValidation_Rollback verifies that when a move+update
// fails date validation at the new position, the move is fully rolled back.
func TestUpdateJump_MoveFailsDateValidation_Rollback(t *testing.T) {
	db := testDB(t)
	// Create 3 jumps: #1 (Mar 1), #2 (Mar 5), #3 (Mar 10)
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 1), Dropzone: "DZ", JumpType: common.JumpTypeFF})
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 5), Dropzone: "DZ", JumpType: common.JumpTypeFF})
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 10), Dropzone: "DZ", JumpType: common.JumpTypeFF})

	// Try to move jump #1 (Mar 1) to position #3, keeping date Mar 1.
	// At position #3: prev=#2 (Mar 5) → Mar 1 < Mar 5 → date validation FAILS.
	// The move should be rolled back — jump should still be at position #1.
	j, _ := db.GetJumpByNumber(1, 1)
	body := jsonBody(t, map[string]any{
		"date":     "2025-03-01",
		"dropzone": "DZ",
		"jumpType": "FF",
		"number":   3,
	})
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/jumps/%d", j.ID), body)
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", j.ID)})
	rr := httptest.NewRecorder()
	handlers.UpdateJump(db)(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rr.Code, rr.Body.String())
	}

	// Verify the jump is still at position #1 (rollback worked)
	j1, err := db.GetJumpByNumber(1, 1)
	if err != nil {
		t.Fatalf("jump #1 not found after rollback: %v", err)
	}
	if j1.ID != j.ID {
		t.Errorf("jump #1 should still be the original jump (ID %d), got ID %d", j.ID, j1.ID)
	}
	if j1.Date.DayString() != "2025-03-01" {
		t.Errorf("jump #1 date should be 2025-03-01, got %s", j1.Date.DayString())
	}

	// Verify all 3 jumps are in their original positions
	j2, _ := db.GetJumpByNumber(1, 2)
	if j2.Date.DayString() != "2025-03-05" {
		t.Errorf("jump #2 should still be Mar 5, got %s", j2.Date.DayString())
	}
	j3, _ := db.GetJumpByNumber(1, 3)
	if j3.Date.DayString() != "2025-03-10" {
		t.Errorf("jump #3 should still be Mar 10, got %s", j3.Date.DayString())
	}
}

// TestUpdateJump_MoveAndDateChange_Atomic verifies that a combined move + date
// change succeeds atomically when the new date is valid at the new position.
func TestUpdateJump_MoveAndDateChange_Atomic(t *testing.T) {
	db := testDB(t)
	// Create 3 jumps: #1 (Mar 1), #2 (Mar 5), #3 (Mar 10)
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 1), Dropzone: "DZ", JumpType: common.JumpTypeFF})
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 5), Dropzone: "DZ", JumpType: common.JumpTypeFF})
	db.CreateJump(&common.Jump{UserID: 1, Date: common.NewDateOnly(2025, time.March, 10), Dropzone: "DZ", JumpType: common.JumpTypeFF})

	// Move jump #3 (Mar 10) to position #1 with date Mar 1.
	// At position #1: no prev, next=#2 (was #1, now shifted to #2, date Mar 1) → Mar 1 <= Mar 1 → VALID
	j, _ := db.GetJumpByNumber(1, 3)
	body := jsonBody(t, map[string]any{
		"date":     "2025-01-01",
		"dropzone": "DZ",
		"jumpType": "FF",
		"number":   1,
	})
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/jumps/%d", j.ID), body)
	req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", j.ID)})
	rr := httptest.NewRecorder()
	handlers.UpdateJump(db)(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var result common.Jump
	json.NewDecoder(rr.Body).Decode(&result)
	if result.Number != 1 {
		t.Errorf("expected number 1, got %d", result.Number)
	}
	if result.Date.DayString() != "2025-01-01" {
		t.Errorf("expected date 2025-01-01, got %s", result.Date.DayString())
	}

	// Verify the other jumps shifted correctly
	j2, _ := db.GetJumpByNumber(1, 2)
	if j2.Date.DayString() != "2025-03-01" {
		t.Errorf("jump #2 should be Mar 1 (was original #1), got %s", j2.Date.DayString())
	}
	j3, _ := db.GetJumpByNumber(1, 3)
	if j3.Date.DayString() != "2025-03-05" {
		t.Errorf("jump #3 should be Mar 5 (was original #2), got %s", j3.Date.DayString())
	}
}
