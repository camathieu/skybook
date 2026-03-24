---
ticket: "007"
epic: code-quality
milestone: v1.2
title: Metadata Date Validation Tests
status: planned
priority: medium
estimate: S
---

# Metadata Date Validation Tests

`validateDateOrder()` in `metadata/jump.go` is tested indirectly through handler tests, but has no direct metadata-layer tests. This means refactoring the handler could silently drop test coverage for the core invariant.

## Acceptance Criteria

- [ ] Add tests to `server/metadata/backend_test.go` that exercise `validateDateOrder` via the public API:
  - `TestCreateJump_DateOrder_AppendValid` — append with date ≥ last
  - `TestCreateJump_DateOrder_AppendInvalid` — append with date < last → error
  - `TestInsertJumpAt_DateOrder_Valid` — insert between two jumps with valid date
  - `TestInsertJumpAt_DateOrder_InvalidBefore` — insert date < previous → error
  - `TestInsertJumpAt_DateOrder_InvalidAfter` — insert date > next → error
  - `TestUpdateJump_DateOrder_Valid` — change date within bounds
  - `TestUpdateJump_DateOrder_Invalid` — change date out of bounds → error
  - `TestCreateJump_SameDay_Valid` — multiple jumps on same day
  - `TestFirstJump_NoValidation` — first jump in empty logbook always passes
- [ ] Tests verify that `DateOrderError` is returned (using `errors.As`)
- [ ] Tests are at the metadata layer (call `b.CreateJump()`, not HTTP handlers)
