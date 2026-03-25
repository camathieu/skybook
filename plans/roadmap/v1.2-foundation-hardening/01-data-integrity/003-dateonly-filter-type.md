---
ticket: "003"
epic: data-integrity
milestone: v1.2
title: DateOnly Filter Type Consistency
status: done
priority: high
estimate: S
---

# DateOnly Filter Type Consistency

`common/date.go` is 124 lines implementing 5 interfaces (`json.Marshaler`, `json.Unmarshaler`, `driver.Valuer`, `sql.Scanner`, plus `DateOrderError`). Every JSON request/response and every database read/write flows through this type. There are **zero unit tests** for it.

## Acceptance Criteria

- [x] Change `JumpFilters.DateFrom` and `DateTo` from `*time.Time` to `*common.DateOnly`
- [x] Update `handlers/jump.go` ListJumps to parse date filter params using `common.NewDateOnly()` instead of `time.Parse()`
- [x] Verify SQL filter queries still work correctly (GORM uses `DateOnly.Value()` which returns `time.Time`)
- [x] Existing filter tests pass
- [x] Add test: filter by `date_from=2025-03-05` returns only jumps on or after that date

## Done

- Changed `JumpFilters.DateFrom` and `DateTo` from `*time.Time` to `*common.DateOnly` in `server/metadata/jump.go`; removed now-unused `time` import
- Updated `handlers/jump.go` to parse date params via `time.Parse(common.DateLayout, ...)` + `common.NewDateOnly()`; dereferenced filter values in GORM WHERE clauses
- Added 3 new tests in `handlers/jump_test.go`: `TestListJumps_DateFromFilter`, `TestListJumps_DateToFilter`, `TestListJumps_InvalidDateFrom`
- All 6 test packages pass; lint and build clean
