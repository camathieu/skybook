package common

import (
	"database/sql/driver"
	"errors"
	"testing"
	"time"
)

// --- Helpers ---

// march15 is a fixed, known DateOnly used across tests.
var march15 = NewDateOnly(2025, time.March, 15)

// --- Constructor & Basic Methods ---

func TestNewDateOnly(t *testing.T) {
	d := NewDateOnly(2025, time.March, 15)
	if d.Year() != 2025 || d.Month() != time.March || d.Day() != 15 {
		t.Errorf("NewDateOnly: got %v, want 2025-03-15", d)
	}
	if d.Hour() != 0 || d.Minute() != 0 || d.Second() != 0 || d.Nanosecond() != 0 {
		t.Errorf("NewDateOnly: time component should be midnight, got %v", d.Time)
	}
	if d.Location() != time.UTC {
		t.Errorf("NewDateOnly: expected UTC, got %v", d.Location())
	}
}

func TestToday(t *testing.T) {
	now := time.Now().UTC()
	today := Today()
	if today.Year() != now.Year() || today.Month() != now.Month() || today.Day() != now.Day() {
		t.Errorf("Today(): got %v, want %04d-%02d-%02d", today.DayString(), now.Year(), now.Month(), now.Day())
	}
}

func TestTruncateToDay(t *testing.T) {
	cases := []struct {
		name string
		in   DateOnly
	}{
		{"midnight is unchanged", NewDateOnly(2025, time.March, 15)},
		{"with time component", DateOnly{Time: time.Date(2025, time.March, 15, 14, 30, 45, 123456789, time.UTC)}},
		{"with non-UTC timezone", DateOnly{Time: time.Date(2025, time.March, 15, 23, 0, 0, 0, time.FixedZone("EST", -5*3600))}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.in.TruncateToDay()
			if got.Year() != tc.in.Year() || got.Month() != tc.in.Month() || got.Day() != tc.in.Day() {
				t.Errorf("TruncateToDay: date mismatch: got %v from %v", got.DayString(), tc.in.DayString())
			}
			if got.Hour() != 0 || got.Minute() != 0 || got.Second() != 0 || got.Nanosecond() != 0 {
				t.Errorf("TruncateToDay: time component should be midnight, got %v", got.Time)
			}
			if got.Location() != time.UTC {
				t.Errorf("TruncateToDay: expected UTC, got %v", got.Location())
			}
		})
	}
}

func TestAddDays(t *testing.T) {
	cases := []struct {
		name string
		base DateOnly
		n    int
		want string
	}{
		{"add positive", NewDateOnly(2025, time.March, 15), 5, "2025-03-20"},
		{"add negative", NewDateOnly(2025, time.March, 15), -5, "2025-03-10"},
		{"add zero", NewDateOnly(2025, time.March, 15), 0, "2025-03-15"},
		{"month boundary forward (Jan 31 + 1)", NewDateOnly(2025, time.January, 31), 1, "2025-02-01"},
		{"month boundary backward (Mar 1 - 1)", NewDateOnly(2025, time.March, 1), -1, "2025-02-28"},
		{"year boundary forward (Dec 31 + 1)", NewDateOnly(2025, time.December, 31), 1, "2026-01-01"},
		{"year boundary backward (Jan 1 - 1)", NewDateOnly(2025, time.January, 1), -1, "2024-12-31"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.base.AddDays(tc.n)
			if got.DayString() != tc.want {
				t.Errorf("AddDays(%d): got %q, want %q", tc.n, got.DayString(), tc.want)
			}
			// Verify time is still midnight UTC.
			if got.Hour() != 0 || got.Minute() != 0 || got.Second() != 0 {
				t.Errorf("AddDays: result should be midnight UTC")
			}
		})
	}
}

func TestSameDay(t *testing.T) {
	cases := []struct {
		name string
		a, b DateOnly
		want bool
	}{
		{
			"same day (identical)",
			NewDateOnly(2025, time.March, 15),
			NewDateOnly(2025, time.March, 15),
			true,
		},
		{
			"same day different time",
			DateOnly{Time: time.Date(2025, time.March, 15, 0, 0, 0, 0, time.UTC)},
			DateOnly{Time: time.Date(2025, time.March, 15, 23, 59, 59, 0, time.UTC)},
			true,
		},
		{
			"different day same month",
			NewDateOnly(2025, time.March, 15),
			NewDateOnly(2025, time.March, 16),
			false,
		},
		{
			"different month same day",
			NewDateOnly(2025, time.March, 15),
			NewDateOnly(2025, time.April, 15),
			false,
		},
		{
			"different year same month/day",
			NewDateOnly(2025, time.March, 15),
			NewDateOnly(2024, time.March, 15),
			false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.a.SameDay(tc.b); got != tc.want {
				t.Errorf("SameDay: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDayString(t *testing.T) {
	cases := []struct {
		d    DateOnly
		want string
	}{
		{NewDateOnly(2025, time.March, 8), "2025-03-08"},
		{NewDateOnly(2025, time.December, 31), "2025-12-31"},
		{NewDateOnly(2025, time.January, 1), "2025-01-01"},
	}
	for _, tc := range cases {
		got := tc.d.DayString()
		if got != tc.want {
			t.Errorf("DayString: got %q, want %q", got, tc.want)
		}
	}
}

func TestIsZero(t *testing.T) {
	var zero DateOnly
	if !zero.IsZero() {
		t.Error("IsZero: zero value should return true")
	}
	if march15.IsZero() {
		t.Error("IsZero: non-zero value should return false")
	}
}

// --- JSON Marshaling ---

func TestMarshalJSON(t *testing.T) {
	// Standard date.
	data, err := march15.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON: unexpected error: %v", err)
	}
	got := string(data)
	want := `"2025-03-15"`
	if got != want {
		t.Errorf("MarshalJSON: got %s, want %s", got, want)
	}

	// Zero-padded month and day.
	d := NewDateOnly(2025, time.March, 8)
	data, _ = d.MarshalJSON()
	if string(data) != `"2025-03-08"` {
		t.Errorf("MarshalJSON: zero-padding failed, got %s", string(data))
	}

	// No time component in output.
	withTime := DateOnly{Time: time.Date(2025, time.March, 15, 14, 30, 0, 0, time.UTC)}
	data, _ = withTime.MarshalJSON()
	if string(data) != `"2025-03-15"` {
		t.Errorf("MarshalJSON: should strip time component, got %s", string(data))
	}
}

func TestUnmarshalJSON(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantStr string // expected DayString(), "" means zero value
		wantErr bool
	}{
		{"valid YYYY-MM-DD", `"2025-03-08"`, "2025-03-08", false},
		{"valid RFC3339 strips time", `"2025-03-08T14:30:00Z"`, "2025-03-08", false},
		{"empty string → zero", `""`, "", false},
		{"null → zero", `"null"`, "", false},
		{"invalid format → error", `"08/03/2025"`, "", true},
		{"invalid date → error", `"not-a-date"`, "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var d DateOnly
			err := d.UnmarshalJSON([]byte(tc.input))
			if tc.wantErr {
				if err == nil {
					t.Errorf("UnmarshalJSON(%s): expected error, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("UnmarshalJSON(%s): unexpected error: %v", tc.input, err)
			}
			if tc.wantStr == "" {
				if !d.IsZero() {
					t.Errorf("UnmarshalJSON(%s): expected zero value, got %v", tc.input, d.DayString())
				}
			} else {
				if d.DayString() != tc.wantStr {
					t.Errorf("UnmarshalJSON(%s): got %q, want %q", tc.input, d.DayString(), tc.wantStr)
				}
			}
		})
	}
}

// --- driver.Valuer ---

func TestValue(t *testing.T) {
	val, err := march15.Value()
	if err != nil {
		t.Fatalf("Value: unexpected error: %v", err)
	}
	got, ok := val.(time.Time)
	if !ok {
		t.Fatalf("Value: expected time.Time, got %T", val)
	}
	if !got.Equal(march15.Time) {
		t.Errorf("Value: got %v, want %v", got, march15.Time)
	}
	// Verify it satisfies the driver.Valuer interface at compile time.
	var _ driver.Valuer = march15
}

// --- sql.Scanner ---

func TestScan(t *testing.T) {
	t.Run("time.Time input", func(t *testing.T) {
		var d DateOnly
		input := time.Date(2025, time.March, 15, 0, 0, 0, 0, time.UTC)
		if err := d.Scan(input); err != nil {
			t.Fatalf("Scan(time.Time): unexpected error: %v", err)
		}
		if d.DayString() != "2025-03-15" {
			t.Errorf("Scan(time.Time): got %q, want %q", d.DayString(), "2025-03-15")
		}
	})

	t.Run("RFC3339 string", func(t *testing.T) {
		var d DateOnly
		if err := d.Scan("2025-03-15T00:00:00Z"); err != nil {
			t.Fatalf("Scan(RFC3339 string): unexpected error: %v", err)
		}
		if d.DayString() != "2025-03-15" {
			t.Errorf("Scan(RFC3339 string): got %q, want %q", d.DayString(), "2025-03-15")
		}
	})

	t.Run("YYYY-MM-DD string", func(t *testing.T) {
		var d DateOnly
		if err := d.Scan("2025-03-15"); err != nil {
			t.Fatalf("Scan(YYYY-MM-DD string): unexpected error: %v", err)
		}
		if d.DayString() != "2025-03-15" {
			t.Errorf("Scan(YYYY-MM-DD string): got %q, want %q", d.DayString(), "2025-03-15")
		}
	})

	t.Run("nil → zero value", func(t *testing.T) {
		var d DateOnly
		if err := d.Scan(nil); err != nil {
			t.Fatalf("Scan(nil): unexpected error: %v", err)
		}
		if !d.IsZero() {
			t.Errorf("Scan(nil): expected zero value, got %v", d.DayString())
		}
	})

	t.Run("unsupported type → error", func(t *testing.T) {
		var d DateOnly
		if err := d.Scan(42); err == nil {
			t.Error("Scan(int): expected error, got nil")
		}
	})

	t.Run("invalid string → error", func(t *testing.T) {
		var d DateOnly
		if err := d.Scan("not-a-date"); err == nil {
			t.Error("Scan(invalid string): expected error, got nil")
		}
	})
}

// --- Round-trip tests ---

func TestRoundTripJSON(t *testing.T) {
	original := NewDateOnly(2025, time.March, 15)

	data, err := original.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}

	var restored DateOnly
	if err := restored.UnmarshalJSON(data); err != nil {
		t.Fatalf("UnmarshalJSON: %v", err)
	}

	if !original.SameDay(restored) {
		t.Errorf("JSON round-trip: got %q, want %q", restored.DayString(), original.DayString())
	}
}

func TestRoundTripDB(t *testing.T) {
	original := NewDateOnly(2025, time.March, 15)

	val, err := original.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}

	var restored DateOnly
	if err := restored.Scan(val); err != nil {
		t.Fatalf("Scan: %v", err)
	}

	if !original.SameDay(restored) {
		t.Errorf("DB round-trip: got %q, want %q", restored.DayString(), original.DayString())
	}
}

// --- DateOrderError ---

func TestDateOrderError(t *testing.T) {
	msg := "date 2025-03-15 is before jump #5 (2025-03-20)"
	err := &DateOrderError{Message: msg}

	// Implements error interface.
	if err.Error() != msg {
		t.Errorf("DateOrderError.Error(): got %q, want %q", err.Error(), msg)
	}

	// errors.As works.
	wrappedErr := err // already a *DateOrderError; no wrapping needed for As
	var target *DateOrderError
	if !errors.As(wrappedErr, &target) {
		t.Error("errors.As: failed to unwrap DateOrderError")
	}
	if target.Message != msg {
		t.Errorf("errors.As: Message = %q, want %q", target.Message, msg)
	}
}
