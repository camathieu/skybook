package metadata

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/root-gg/skybook/common"
)

// ErrUnsupportedAutocompleteField is returned when an unsupported field is passed to GetJumpAutocomplete.
var ErrUnsupportedAutocompleteField = errors.New("unsupported autocomplete field")

// allowedSortFields is the whitelist of columns allowed for sorting in GetJumps.
var allowedSortFields = map[string]bool{
	"number":   true,
	"date":     true,
	"dropzone": true,
	"altitude": true,
}

// JumpFilters holds optional filter parameters for listing jumps.
type JumpFilters struct {
	Q           string
	DateFrom    *time.Time
	DateTo      *time.Time
	Dropzone    string
	Aircraft    string
	JumpType    string
	AltitudeMin *uint
	AltitudeMax *uint
	Cutaway     *bool
	Night       *bool
	LO          string
}

// CreateJump appends a new jump at the end of the user's logbook.
// The Number is automatically set to MAX(Number)+1.
func (b *Backend) CreateJump(jump *common.Jump) error {
	// Truncate date to midnight before storing
	jump.Date = jump.Date.TruncateToDay()

	return b.db.Transaction(func(tx *gorm.DB) error {
		var maxNumber uint
		if err := tx.Model(&common.Jump{}).
			Where("user_id = ?", jump.UserID).
			Select("COALESCE(MAX(number), 0)").
			Scan(&maxNumber).Error; err != nil {
			return fmt.Errorf("compute max number: %w", err)
		}

		// Validate date ordering: new jump must be >= previous jump's date
		if maxNumber > 0 {
			if err := validateDateOrder(tx, jump.UserID, jump.Date, maxNumber+1, 0); err != nil {
				return err
			}
		}

		jump.Number = maxNumber + 1
		return tx.Create(jump).Error
	})
}

// InsertJumpAt inserts a jump at a specific position in the user's logbook.
// All jumps with Number >= pos are shifted up by 1.
func (b *Backend) InsertJumpAt(jump *common.Jump, pos uint) error {
	if pos < 1 {
		return fmt.Errorf("position must be >= 1, got %d", pos)
	}

	// Truncate date to midnight before storing
	jump.Date = jump.Date.TruncateToDay()

	return b.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		tx.Model(&common.Jump{}).Where("user_id = ?", jump.UserID).Count(&count)

		if pos > uint(count)+1 {
			return fmt.Errorf("position %d out of range (max %d)", pos, count+1)
		}

		// Validate date ordering at insert position
		if err := validateDateOrder(tx, jump.UserID, jump.Date, pos, 0); err != nil {
			return err
		}

		// Shift existing jumps up. SQLite doesn't support UPDATE ... ORDER BY,
		// so we iterate in descending order to avoid unique constraint violations.
		if count > 0 && pos <= uint(count) {
			var jumpsToShift []*common.Jump
			if err := tx.Where("user_id = ? AND number >= ?", jump.UserID, pos).
				Order("number DESC").
				Find(&jumpsToShift).Error; err != nil {
				return err
			}
			for _, j := range jumpsToShift {
				if err := tx.Model(j).Update("number", j.Number+1).Error; err != nil {
					return err
				}
			}
		}

		jump.Number = pos
		return tx.Create(jump).Error
	})
}

// DeleteJump removes a jump and renumbers subsequent jumps to maintain contiguity.
func (b *Backend) DeleteJump(jump *common.Jump) error {
	return b.db.Transaction(func(tx *gorm.DB) error {
		number := jump.Number
		userID := jump.UserID

		if err := tx.Delete(jump).Error; err != nil {
			return err
		}

		// Shift subsequent jumps down in ascending order to avoid unique constraint.
		var jumpsToShift []*common.Jump
		if err := tx.Where("user_id = ? AND number > ?", userID, number).
			Order("number ASC").
			Find(&jumpsToShift).Error; err != nil {
			return err
		}
		for _, j := range jumpsToShift {
			if err := tx.Model(j).Update("number", j.Number-1).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// MoveJump repositions a jump to a new number within the user's logbook.
// It parks the jump at a large sentinel value, shifts the intermediate range,
// then sets the jump to the new position — all within a single transaction.
func (b *Backend) MoveJump(jump *common.Jump, newNumber uint) error {
	if newNumber < 1 {
		return fmt.Errorf("number must be >= 1")
	}

	return b.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		tx.Model(&common.Jump{}).Where("user_id = ?", jump.UserID).Count(&count)
		if newNumber > uint(count) {
			return fmt.Errorf("number %d out of range (max %d)", newNumber, count)
		}

		oldNumber := jump.Number
		userID := jump.UserID

		if newNumber == oldNumber {
			return nil
		}

		// Park at a high sentinel to free the unique slot while we shift.
		// We use raw Exec to bypass GORM's zero-value filtering on uint fields.
		const jumpNumberSentinel uint = 999_999_999
		if err := tx.Exec("UPDATE jumps SET number = ? WHERE id = ?", jumpNumberSentinel, jump.ID).Error; err != nil {
			return fmt.Errorf("park jump: %w", err)
		}

		if newNumber < oldNumber {
			// Moving up: shift rows in [newNumber, oldNumber-1] from position X to X+1.
			// Must process in DESC order (highest number first) to avoid unique constraint.
			var jumpsToShift []*common.Jump
			if err := tx.Where("user_id = ? AND number >= ? AND number < ?", userID, newNumber, oldNumber).
				Order("number DESC").
				Find(&jumpsToShift).Error; err != nil {
				return err
			}
			for _, j := range jumpsToShift {
				if err := tx.Model(j).Update("number", j.Number+1).Error; err != nil {
					return err
				}
			}
		} else {
			// Moving down: shift rows in [oldNumber+1, newNumber] from position X to X-1.
			// Must process in ASC order (lowest number first) to avoid unique constraint.
			var jumpsToShift []*common.Jump
			if err := tx.Where("user_id = ? AND number > ? AND number <= ?", userID, oldNumber, newNumber).
				Order("number ASC").
				Find(&jumpsToShift).Error; err != nil {
				return err
			}
			for _, j := range jumpsToShift {
				if err := tx.Model(j).Update("number", j.Number-1).Error; err != nil {
					return err
				}
			}
		}

		// Place the jump at its final position
		if err := tx.Exec("UPDATE jumps SET number = ? WHERE id = ?", newNumber, jump.ID).Error; err != nil {
			return fmt.Errorf("place jump: %w", err)
		}
		jump.Number = newNumber
		return nil
	})
}

// validateDateOrder checks that date(N) ≤ date(N+1) at the given position.
// skipID is the ID of a jump being moved — it is excluded from neighbor lookup
// to avoid comparing a jump against itself.
func validateDateOrder(tx *gorm.DB, userID uint, date common.DateOnly, atNumber uint, skipID uint) error {
	date = date.TruncateToDay()

	// Check previous jump (N-1)
	if atNumber > 1 {
		var prev common.Jump
		if err := tx.Where("user_id = ? AND number = ? AND id != ?", userID, atNumber-1, skipID).
			First(&prev).Error; err == nil {
			if date.TruncateToDay().Time.Before(prev.Date.TruncateToDay().Time) {
				return &common.DateOrderError{Message: fmt.Sprintf("date %s is before jump #%d (%s)", date.DayString(), prev.Number, prev.Date.DayString())}
			}
		}
	}

	// Check next jump (N or N+1 depending on context)
	var next common.Jump
	if err := tx.Where("user_id = ? AND number = ? AND id != ?", userID, atNumber, skipID).
		First(&next).Error; err == nil {
		// There's a jump at the target position (will be shifted up)
		if date.TruncateToDay().Time.After(next.Date.TruncateToDay().Time) {
			return &common.DateOrderError{Message: fmt.Sprintf("date %s is after jump #%d (%s)", date.DayString(), next.Number, next.Date.DayString())}
		}
	} else {
		// No jump at target — check N+1 (for move/update scenarios)
		if err := tx.Where("user_id = ? AND number = ? AND id != ?", userID, atNumber+1, skipID).
			First(&next).Error; err == nil {
			if date.TruncateToDay().Time.After(next.Date.TruncateToDay().Time) {
				return &common.DateOrderError{Message: fmt.Sprintf("date %s is after jump #%d (%s)", date.DayString(), next.Number, next.Date.DayString())}
			}
		}
	}

	return nil
}

// GetJump retrieves a single jump by ID for a given user.
func (b *Backend) GetJump(userID, jumpID uint) (*common.Jump, error) {
	jump := &common.Jump{}
	err := b.db.Where("id = ? AND user_id = ?", jumpID, userID).First(jump).Error
	if err != nil {
		return nil, err
	}
	return jump, nil
}

// GetJumpByNumber retrieves a single jump by number for a given user.
func (b *Backend) GetJumpByNumber(userID, number uint) (*common.Jump, error) {
	jump := &common.Jump{}
	err := b.db.Where("user_id = ? AND number = ?", userID, number).First(jump).Error
	if err != nil {
		return nil, err
	}
	return jump, nil
}

// UpdateJump persists changes to a jump's fields (excluding Number — use MoveJump).
// Validates date ordering against neighbors before saving.
func (b *Backend) UpdateJump(jump *common.Jump) error {
	jump.Date = jump.Date.TruncateToDay()

	return b.db.Transaction(func(tx *gorm.DB) error {
		if err := validateDateOrder(tx, jump.UserID, jump.Date, jump.Number, jump.ID); err != nil {
			return err
		}
		return tx.Omit(clause.Associations).Save(jump).Error
	})
}

// GetJumps retrieves a paginated, filtered list of jumps for a user.
func (b *Backend) GetJumps(userID uint, offset, limit int, sortBy, order string, filters JumpFilters) ([]*common.Jump, int64, error) {
	var jumps []*common.Jump
	var total int64

	q := b.db.Model(&common.Jump{}).Where("user_id = ?", userID)

	if filters.Q != "" {
		like := "%" + filters.Q + "%"
		q = q.Where("description LIKE ? OR dropzone LIKE ? OR event LIKE ? OR lo LIKE ?", like, like, like, like)
	}
	if filters.DateFrom != nil {
		q = q.Where("date >= ?", filters.DateFrom)
	}
	if filters.DateTo != nil {
		q = q.Where("date <= ?", filters.DateTo)
	}
	if filters.Dropzone != "" {
		q = q.Where("dropzone = ?", filters.Dropzone)
	}
	if filters.Aircraft != "" {
		q = q.Where("aircraft = ?", filters.Aircraft)
	}
	if filters.JumpType != "" {
		q = q.Where("jump_type = ?", filters.JumpType)
	}
	if filters.AltitudeMin != nil {
		q = q.Where("altitude >= ?", *filters.AltitudeMin)
	}
	if filters.AltitudeMax != nil {
		q = q.Where("altitude <= ?", *filters.AltitudeMax)
	}
	if filters.Cutaway != nil {
		q = q.Where("cut_away = ?", *filters.Cutaway)
	}
	if filters.Night != nil {
		q = q.Where("night_jump = ?", *filters.Night)
	}
	if filters.LO != "" {
		q = q.Where("lo = ?", filters.LO)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sortBy == "" {
		sortBy = "number"
	} else if !allowedSortFields[sortBy] {
		return nil, 0, fmt.Errorf("invalid sort field: %s", sortBy)
	}
	if order == "" {
		order = "desc"
	} else if order != "asc" && order != "desc" {
		return nil, 0, fmt.Errorf("invalid order: %s", order)
	}
	orderClause := fmt.Sprintf("%s %s", sortBy, order)

	err := q.Order(orderClause).Offset(offset).Limit(limit).Find(&jumps).Error
	return jumps, total, err
}

// CountJumps returns the total number of jumps for a user.
func (b *Backend) CountJumps(userID uint) (int64, error) {
	var count int64
	err := b.db.Model(&common.Jump{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// allowedAutocompleteFields maps allowed field names to their SQL column names.
// Note: only add fields that correspond to actual columns in the jumps table.
var allowedAutocompleteFields = map[string]string{
	"dropzone":  "dropzone",
	"aircraft":  "aircraft",
	"jump_type": "jump_type",
	"lo":        "lo",
	"event":     "event",
}

// GetJumpAutocomplete returns distinct non-empty values for a given field.
// sortBy: "alpha" → alphabetical (col ASC); anything else → recency (MAX(date) DESC).
// If prefix is empty, all distinct values are returned (powers on-focus suggestions).
// Only fields in allowedAutocompleteFields are supported.
func (b *Backend) GetJumpAutocomplete(userID uint, field, prefix, sortBy string, limit int) ([]string, error) {
	col, ok := allowedAutocompleteFields[field]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedAutocompleteField, field)
	}

	var orderClause string
	if sortBy == "alpha" {
		orderClause = col + " ASC"
	} else {
		orderClause = "MAX(date) DESC, " + col + " ASC"
	}

	var results []string
	q := b.db.Model(&common.Jump{}).
		Select(col).
		Where("user_id = ? AND "+col+" != ''", userID).
		Group(col).
		Order(orderClause).
		Limit(limit)

	if prefix != "" {
		// Use LOWER() on both sides for true case-insensitive matching.
		q = q.Where("LOWER("+col+") LIKE ?", strings.ToLower(prefix)+"%")
	}

	if err := q.Pluck(col, &results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
