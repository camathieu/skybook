---
ticket: "003"
epic: data-integrity
milestone: v1.2
title: DateOnly Filter Type Consistency
status: planned
priority: high
estimate: S
---

# DateOnly Filter Type Consistency

`JumpFilters.DateFrom` and `DateTo` are `*time.Time` while `Jump.Date` is `DateOnly`. This type mismatch means filter comparisons go through different serialization paths and would break if anyone tried to compare them in Go code.

## Acceptance Criteria

- [ ] Change `JumpFilters.DateFrom` and `DateTo` from `*time.Time` to `*common.DateOnly`
- [ ] Update `handlers/jump.go` ListJumps to parse date filter params using `common.NewDateOnly()` instead of `time.Parse()`
- [ ] Verify SQL filter queries still work correctly (GORM uses `DateOnly.Value()` which returns `time.Time`)
- [ ] Existing filter tests pass
- [ ] Add test: filter by `date_from=2025-03-05` returns only jumps on or after that date
