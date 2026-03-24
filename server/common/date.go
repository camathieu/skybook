package common

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// DateLayout is the canonical date format used in API input/output.
const DateLayout = "2006-01-02"

// DateOrderError is returned when a jump's date violates the date(N) ≤ date(N+1) invariant.
type DateOrderError struct {
	Message string
}

func (e *DateOrderError) Error() string { return e.Message }

// DateOnly wraps time.Time but serializes to/from "YYYY-MM-DD" in JSON.
// Stored as midnight UTC in the database. The time component is never
// exposed via the API.
type DateOnly struct {
	time.Time
}

// NewDateOnly creates a DateOnly from year, month, day at midnight UTC.
func NewDateOnly(year int, month time.Month, day int) DateOnly {
	return DateOnly{Time: time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

// Today returns today's date at midnight UTC.
func Today() DateOnly {
	now := time.Now().UTC()
	return NewDateOnly(now.Year(), now.Month(), now.Day())
}

// TruncateToDay returns a new DateOnly with the time set to midnight UTC,
// preserving year/month/day.
func (d DateOnly) TruncateToDay() DateOnly {
	return NewDateOnly(d.Year(), d.Month(), d.Day())
}

// AddDays returns a new DateOnly shifted by n calendar days.
func (d DateOnly) AddDays(n int) DateOnly {
	t := d.Time.AddDate(0, 0, n)
	return NewDateOnly(t.Year(), t.Month(), t.Day())
}

// SameDay returns true if both dates fall on the same calendar day (UTC).
func (d DateOnly) SameDay(other DateOnly) bool {
	return d.Year() == other.Year() && d.Month() == other.Month() && d.Day() == other.Day()
}

// DayString returns the date formatted as "YYYY-MM-DD".
func (d DateOnly) DayString() string {
	return d.Format(DateLayout)
}

// IsZero returns true if the underlying time is zero.
func (d DateOnly) IsZero() bool {
	return d.Time.IsZero()
}

// MarshalJSON serializes the date as "YYYY-MM-DD".
func (d DateOnly) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.DayString() + `"`), nil
}

// UnmarshalJSON accepts both "YYYY-MM-DD" and RFC3339 formats.
// Any time component is stripped — only the date is preserved.
func (d *DateOnly) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		d.Time = time.Time{}
		return nil
	}

	// Try date-only first (preferred)
	t, err := time.Parse(DateLayout, s)
	if err == nil {
		d.Time = t
		return nil
	}

	// Fall back to RFC3339 (strip time component)
	t, err = time.Parse(time.RFC3339, s)
	if err != nil {
		return fmt.Errorf("date must be YYYY-MM-DD or RFC3339, got %q", s)
	}
	d.Time = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return nil
}

// Value implements driver.Valuer for GORM/SQL.
func (d DateOnly) Value() (driver.Value, error) {
	return d.Time, nil
}

// Scan implements sql.Scanner for GORM/SQL.
func (d *DateOnly) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			t, err = time.Parse(DateLayout, v)
		}
		if err != nil {
			return fmt.Errorf("cannot scan %q into DateOnly", v)
		}
		d.Time = t
		return nil
	default:
		return fmt.Errorf("cannot scan %T into DateOnly", value)
	}
}
