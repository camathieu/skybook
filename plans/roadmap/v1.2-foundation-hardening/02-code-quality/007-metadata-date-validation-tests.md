---
ticket: "007"
epic: code-quality
milestone: v1.2
title: Metadata Date Validation Tests
status: done
priority: medium
estimate: S
---

# Metadata Date Validation Tests

`validateDateOrder()` in `metadata/jump.go` is tested indirectly through handler tests, but has no direct metadata-layer tests. This means refactoring the handler could silently drop test coverage for the core invariant.

## Acceptance Criteria

- [x] Add tests to `server/metadata/backend_test.go` that exercise `validateDateOrder` via the public API:
  - [x] `TestCreateJump_DateOrder_AppendValid` — append with date ≥ last
  - [x] `TestCreateJump_DateOrder_AppendInvalid` — append with date < last → error
  - [x] `TestInsertJumpAt_DateOrder_Valid` — insert between two jumps with valid date
  - [x] `TestInsertJumpAt_DateOrder_InvalidBefore` — insert date < previous → error
  - [x] `TestInsertJumpAt_DateOrder_InvalidAfter` — insert date > next → error
  - [x] `TestUpdateJump_DateOrder_Valid` — change date within bounds
  - [x] `TestUpdateJump_DateOrder_Invalid` — change date out of bounds → error
  - [x] `TestCreateJump_SameDay_Valid` — multiple jumps on same day
  - [x] `TestFirstJump_NoValidation` — first jump in empty logbook always passes
- [x] Tests verify that `DateOrderError` is returned (using `errors.As`)
- [x] Tests are at the metadata layer (call `b.CreateJump()` / `b.InsertJumpAt()`, not HTTP handlers)

## Done

- Added `dateOf(year, month, day)` helper to `server/metadata/backend_test.go`
- Added all 9 test functions exercising `CreateJump`, `InsertJumpAt`, and `UpdateJump`
- All `DateOrderError` assertions use `errors.As` (not string matching)
- Key fix during implementation: `InsertJumpAt_*` tests must call `b.InsertJumpAt(jump, pos)` directly — `CreateJump` always appends and ignores the `Number` field
- All 6 test packages pass
