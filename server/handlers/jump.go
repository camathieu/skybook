package handlers

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/root-gg/skybook/common"
	"github.com/root-gg/skybook/metadata"
)

// anonymousUserID is the hardcoded user ID for v1 single-user mode.
const anonymousUserID uint = 1

// jumpListResponse is the paginated response envelope for GET /api/v1/jumps.
type jumpListResponse struct {
	Items      []*common.Jump `json:"items"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PerPage    int            `json:"perPage"`
	TotalPages int            `json:"totalPages"`
}

// ListJumps handles GET /api/v1/jumps
func ListJumps(db *metadata.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		// --- Pagination ---
		page := parseIntParam(q.Get("page"), 1)
		if page < 1 {
			page = 1
		}
		perPage := parseIntParam(q.Get("per_page"), 25)
		if perPage < 1 {
			perPage = 25
		}
		if perPage > 100 {
			perPage = 100
		}
		offset := (page - 1) * perPage

		// --- Sorting ---
		sortBy := q.Get("sort")
		order := q.Get("order")

		// Whitelist allowed sort columns to prevent SQL injection
		allowedSort := map[string]bool{"number": true, "date": true, "dropzone": true, "altitude": true}
		if sortBy != "" && !allowedSort[sortBy] {
			common.WriteError(w, "invalid sort field", http.StatusBadRequest)
			return
		}
		if order != "" && order != "asc" && order != "desc" {
			common.WriteError(w, "order must be 'asc' or 'desc'", http.StatusBadRequest)
			return
		}

		// --- Filters ---
		filters := metadata.JumpFilters{
			Q:        q.Get("q"),
			Dropzone: q.Get("dropzone"),
			Aircraft: q.Get("aircraft"),
			JumpType: q.Get("jump_type"),
			LO:       q.Get("lo"),
		}

		if df := q.Get("date_from"); df != "" {
			t, err := time.Parse("2006-01-02", df)
			if err != nil {
				common.WriteError(w, "invalid date_from format (use YYYY-MM-DD)", http.StatusBadRequest)
				return
			}
			filters.DateFrom = &t
		}
		if dt := q.Get("date_to"); dt != "" {
			t, err := time.Parse("2006-01-02", dt)
			if err != nil {
				common.WriteError(w, "invalid date_to format (use YYYY-MM-DD)", http.StatusBadRequest)
				return
			}
			filters.DateTo = &t
		}
		if v := q.Get("altitude_min"); v != "" {
			n, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				common.WriteError(w, "invalid altitude_min", http.StatusBadRequest)
				return
			}
			u := uint(n)
			filters.AltitudeMin = &u
		}
		if v := q.Get("altitude_max"); v != "" {
			n, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				common.WriteError(w, "invalid altitude_max", http.StatusBadRequest)
				return
			}
			u := uint(n)
			filters.AltitudeMax = &u
		}
		if v := q.Get("cutaway"); v != "" {
			switch v {
			case "true", "1":
				b := true
				filters.Cutaway = &b
			case "false", "0":
				b := false
				filters.Cutaway = &b
			default:
				common.WriteError(w, "cutaway must be true or false", http.StatusBadRequest)
				return
			}
		}
		if v := q.Get("night"); v != "" {
			switch v {
			case "true", "1":
				b := true
				filters.Night = &b
			case "false", "0":
				b := false
				filters.Night = &b
			default:
				common.WriteError(w, "night must be true or false", http.StatusBadRequest)
				return
			}
		}

		jumps, total, err := db.GetJumps(anonymousUserID, offset, perPage, sortBy, order, filters)
		if err != nil {
			common.WriteError(w, "failed to list jumps", http.StatusInternalServerError)
			return
		}

		totalPages := int(math.Ceil(float64(total) / float64(perPage)))
		if total == 0 {
			totalPages = 0
		}

		// Return empty array instead of null
		if jumps == nil {
			jumps = []*common.Jump{}
		}

		common.WriteJSON(w, jumpListResponse{
			Items:      jumps,
			Total:      total,
			Page:       page,
			PerPage:    perPage,
			TotalPages: totalPages,
		}, http.StatusOK)
	}
}

// CreateJump handles POST /api/v1/jumps
// If the request body includes a non-zero "number" field, it inserts at that position.
// Otherwise, it appends to the end of the logbook.
func CreateJump(db *metadata.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var jump common.Jump
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB max
		if err := json.NewDecoder(r.Body).Decode(&jump); err != nil {
			common.WriteError(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		// Validate required fields
		if jump.Date.IsZero() {
			common.WriteError(w, "date is required", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(jump.Dropzone) == "" {
			common.WriteError(w, "dropzone is required", http.StatusBadRequest)
			return
		}
		if !jump.JumpType.IsValid() {
			common.WriteError(w, "invalid jump_type", http.StatusBadRequest)
			return
		}

		jump.UserID = anonymousUserID

		// Truncate date to day-only (backend controls time component)
		jump.Date = jump.Date.TruncateToDay()

		var err error
		requestedNumber := jump.Number
		jump.Number = 0 // will be assigned by Create/Insert

		if requestedNumber > 0 {
			// Validate insert range
			count, cerr := db.CountJumps(anonymousUserID)
			if cerr != nil {
				common.WriteError(w, "failed to count jumps", http.StatusInternalServerError)
				return
			}
			if requestedNumber > uint(count)+1 {
				common.WriteError(w, "number out of range", http.StatusBadRequest)
				return
			}
			err = db.InsertJumpAt(&jump, requestedNumber)
		} else {
			err = db.CreateJump(&jump)
		}

		if err != nil {
			var dateErr *common.DateOrderError
			if errors.As(err, &dateErr) {
				common.WriteError(w, dateErr.Message, http.StatusBadRequest)
			} else {
				common.WriteError(w, "failed to create jump", http.StatusInternalServerError)
			}
			return
		}

		common.WriteJSON(w, &jump, http.StatusCreated)
	}
}

// GetJump handles GET /api/v1/jumps/{id}
func GetJump(db *metadata.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseIDVar(r, "id")
		if err != nil {
			common.WriteError(w, "invalid jump ID", http.StatusBadRequest)
			return
		}

		jump, err := db.GetJump(anonymousUserID, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				common.WriteError(w, "jump not found", http.StatusNotFound)
				return
			}
			common.WriteError(w, "failed to get jump", http.StatusInternalServerError)
			return
		}

		common.WriteJSON(w, jump, http.StatusOK)
	}
}

// UpdateJump handles PUT /api/v1/jumps/{id}
func UpdateJump(db *metadata.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseIDVar(r, "id")
		if err != nil {
			common.WriteError(w, "invalid jump ID", http.StatusBadRequest)
			return
		}

		// Fetch existing jump to verify ownership and get current number
		existing, err := db.GetJump(anonymousUserID, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				common.WriteError(w, "jump not found", http.StatusNotFound)
				return
			}
			common.WriteError(w, "failed to get jump", http.StatusInternalServerError)
			return
		}

		var body common.Jump
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB max
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			common.WriteError(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		// Validate required fields
		if body.Date.IsZero() {
			common.WriteError(w, "date is required", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(body.Dropzone) == "" {
			common.WriteError(w, "dropzone is required", http.StatusBadRequest)
			return
		}
		if !body.JumpType.IsValid() {
			common.WriteError(w, "invalid jump_type", http.StatusBadRequest)
			return
		}

		// Preserve immutable fields
		body.ID = existing.ID
		body.UserID = existing.UserID
		body.CreatedAt = existing.CreatedAt

		// Truncate date to day-only
		body.Date = body.Date.TruncateToDay()

		requestedNumber := body.Number
		body.Number = existing.Number // will be MoveJump'd if needed

		// Handle number change (reposition) first, before field update
		if requestedNumber != 0 && requestedNumber != existing.Number {
			count, cerr := db.CountJumps(anonymousUserID)
			if cerr != nil {
				common.WriteError(w, "failed to count jumps", http.StatusInternalServerError)
				return
			}
			if requestedNumber < 1 || requestedNumber > uint(count) {
				common.WriteError(w, "number out of range", http.StatusBadRequest)
				return
			}
			if err := db.MoveJump(existing, requestedNumber); err != nil {
				common.WriteError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			body.Number = requestedNumber
		}

		// Validate date ordering at the (possibly new) position.
		// Done via UpdateJump which calls validateDateOrder internally.
		if err := db.UpdateJump(&body); err != nil {
			var dateErr *common.DateOrderError
			if errors.As(err, &dateErr) {
				common.WriteError(w, dateErr.Message, http.StatusBadRequest)
			} else {
				common.WriteError(w, "failed to update jump", http.StatusInternalServerError)
			}
			return
		}

		common.WriteJSON(w, &body, http.StatusOK)
	}
}

// DeleteJump handles DELETE /api/v1/jumps/{id}
func DeleteJump(db *metadata.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseIDVar(r, "id")
		if err != nil {
			common.WriteError(w, "invalid jump ID", http.StatusBadRequest)
			return
		}

		jump, err := db.GetJump(anonymousUserID, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				common.WriteError(w, "jump not found", http.StatusNotFound)
				return
			}
			common.WriteError(w, "failed to get jump", http.StatusInternalServerError)
			return
		}

		if err := db.DeleteJump(jump); err != nil {
			common.WriteError(w, "failed to delete jump", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// Autocomplete handles GET /api/v1/jumps/autocomplete/{field}
// Optional query params:
//   - q: prefix filter
//   - sort: "alpha" for alphabetical; default is recency (MAX(date) DESC)
func Autocomplete(db *metadata.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		field := mux.Vars(r)["field"]
		prefix := r.URL.Query().Get("q")
		sortBy := r.URL.Query().Get("sort")

		results, err := db.GetJumpAutocomplete(anonymousUserID, field, prefix, sortBy, 20)
		if err != nil {
			if errors.Is(err, metadata.ErrUnsupportedAutocompleteField) {
				common.WriteError(w, err.Error(), http.StatusBadRequest)
				return
			}
			common.WriteError(w, "autocomplete failed", http.StatusInternalServerError)
			return
		}

		common.WriteJSON(w, results, http.StatusOK)
	}
}

// --- Helpers ---

func parseIntParam(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}

func parseIDVar(r *http.Request, key string) (uint, error) {
	vars := mux.Vars(r)
	n, err := strconv.ParseUint(vars[key], 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(n), nil
}
