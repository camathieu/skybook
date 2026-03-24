---
ticket: "002"
epic: data-integrity
milestone: v1.2
title: DateOnly Unit Tests
status: planned
priority: critical
estimate: S
---

# DateOnly Unit Tests

`common/date.go` is 124 lines implementing 5 interfaces (`json.Marshaler`, `json.Unmarshaler`, `driver.Valuer`, `sql.Scanner`, plus `DateOrderError`). Every JSON request/response and every database read/write flows through this type. There are **zero unit tests** for it.

## Acceptance Criteria

- [ ] Create `server/common/date_test.go`
- [ ] Test `NewDateOnly` — creates midnight UTC
- [ ] Test `Today` — returns today's date (year/month/day match)
- [ ] Test `TruncateToDay` — strips time component, preserves date
- [ ] Test `AddDays` — positive, negative, month boundary crossing
- [ ] Test `SameDay` — same day returns true, different day returns false, different time same day returns true
- [ ] Test `DayString` — returns `"YYYY-MM-DD"` format
- [ ] Test `IsZero` — zero value returns true, non-zero returns false
- [ ] Test `MarshalJSON` — outputs `"2025-03-08"` (quoted, no time component)
- [ ] Test `UnmarshalJSON` — valid `"2025-03-08"`, valid RFC3339 `"2025-03-08T14:30:00Z"` (time stripped), empty string, `"null"`, invalid format returns error
- [ ] Test `Value` (driver.Valuer) — returns `time.Time`
- [ ] Test `Scan` (sql.Scanner) — `time.Time` input, RFC3339 string, YYYY-MM-DD string, `nil`, unsupported type returns error
- [ ] Test round-trip: `MarshalJSON` → `UnmarshalJSON` preserves date
- [ ] Test round-trip: `Value` → `Scan` preserves date
- [ ] Test `DateOrderError` — implements `error` interface, `errors.As` works
