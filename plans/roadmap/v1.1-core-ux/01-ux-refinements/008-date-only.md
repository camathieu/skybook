---
ticket: "008"
epic: ux-refinements
milestone: v1.1
title: Date-Only Field (Drop Time)
status: done
priority: high
estimate: S
---

# Date-Only Field (Drop Time)

Change the jump date from a `datetime` to a `date`-only field. Skydivers care about the day they jumped, not the exact time.

## Acceptance Criteria

- [x] Backend: Change `Date` field in the `Jump` model from `time.Time` with time precision to date-only storage (`DATE` type in SQLite, or `time.Time` truncated to midnight UTC)
- [x] Backend: API accepts dates in `YYYY-MM-DD` format (e.g., `"2025-03-08"`)
- [x] Backend: API returns dates in `YYYY-MM-DD` format (no time component)
- [x] Frontend: `JumpModal.vue` uses `<input type="date">` instead of `<input type="datetime-local">`
- [x] Frontend: Jump list/cards display dates as `Mar 8, 2025` (no time)
- [x] Migration: Existing datetime values are truncated to date (no data loss risk since time is not meaningful)
- [x] `fakedb` generates date-only values
- [x] All tests updated to use date-only format

## Done

- Created `DateOnly` type in `common/date.go` — wraps `time.Time`, marshals as `YYYY-MM-DD`, accepts both `YYYY-MM-DD` and RFC3339
- Added `DateOrderError` sentinel type for validation error discrimination
- Changed `Jump.Date` from `time.Time` to `DateOnly` in `common/jump.go`
- Frontend: `JumpModal.vue` uses `<input type="date">`, `todayDate()` uses local timezone
- Frontend: `JumpTable.vue` `formatDate` parses date-only strings as local dates (timezone-safe)
- `fakedb.go` uses `DateOnly` + `AddDays()` natively
- Tests: `TestCreateJump_ReturnsDateOnly`, `TestCreateJump_AcceptsRFC3339` verify API format
