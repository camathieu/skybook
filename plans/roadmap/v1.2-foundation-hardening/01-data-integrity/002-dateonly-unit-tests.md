---
ticket: "002"
epic: data-integrity
milestone: v1.2
title: DateOnly Unit Tests
status: done
priority: critical
estimate: S
---

# DateOnly Unit Tests

`common/date.go` is 124 lines implementing 5 interfaces (`json.Marshaler`, `json.Unmarshaler`, `driver.Valuer`, `sql.Scanner`, plus `DateOrderError`). Every JSON request/response and every database read/write flows through this type. There are **zero unit tests** for it.

## Acceptance Criteria

- [x] Create `server/common/date_test.go`
- [x] Test `NewDateOnly` — creates midnight UTC
- [x] Test `Today` — returns today's date (year/month/day match)
- [x] Test `TruncateToDay` — strips time component, preserves date
- [x] Test `AddDays` — positive, negative, month boundary crossing
- [x] Test `SameDay` — same day returns true, different day returns false, different time same day returns true
- [x] Test `DayString` — returns `"YYYY-MM-DD"` format
- [x] Test `IsZero` — zero value returns true, non-zero returns false
- [x] Test `MarshalJSON` — outputs `"2025-03-08"` (quoted, no time component)
- [x] Test `UnmarshalJSON` — valid `"2025-03-08"`, valid RFC3339 `"2025-03-08T14:30:00Z"` (time stripped), empty string, `"null"`, invalid format returns error
- [x] Test `Value` (driver.Valuer) — returns `time.Time`
- [x] Test `Scan` (sql.Scanner) — `time.Time` input, RFC3339 string, YYYY-MM-DD string, `nil`, unsupported type returns error
- [x] Test round-trip: `MarshalJSON` → `UnmarshalJSON` preserves date
- [x] Test round-trip: `Value` → `Scan` preserves date
- [x] Test `DateOrderError` — implements `error` interface, `errors.As` works

## Done

Created `server/common/date_test.go` with 14 test functions covering all `DateOnly` methods and interfaces. All 38 test cases (including subtests) pass. Full `common`, `handlers`, `metadata`, and `middleware` packages pass.
