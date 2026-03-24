---
ticket: "001"
epic: data-integrity
milestone: v1.2
title: Atomic Move+Update Transaction
status: done
priority: critical
estimate: M
---

# Atomic Move+Update Transaction

`UpdateJump` in `handlers/jump.go` calls `MoveJump()` (transaction 1) then `UpdateJump()` (transaction 2) when both number and fields change. If transaction 2 fails (e.g. date validation), the move is already committed — the logbook is now in a state the user never requested.

## Root Cause

The handler orchestrates two independent metadata methods, each with their own `db.Transaction()`. There's no outer transaction wrapping both.

## Acceptance Criteria

- [x] A single metadata method (e.g. `MoveAndUpdateJump(jump, newNumber)`) wraps move + field update + date validation in **one** `db.Transaction()`
- [x] If date validation fails after the move, the entire operation rolls back (move is undone)
- [x] The handler calls this single method instead of `MoveJump()` + `UpdateJump()` separately
- [x] Test: update jump #1 (date=Mar 1) to number=3, date=Mar 1 → should fail date validation AND the jump should still be at position #1
- [x] Test: update jump #2 (date=Mar 5) to number=1, date=Mar 1 → should succeed atomically
- [x] Existing move and update tests still pass

## Done

- `server/metadata/jump.go`: Extracted `moveJumpTx` and `updateJumpTx` private helpers, added `MoveAndUpdateJump` wrapping both in a single `db.Transaction()`
- `server/handlers/jump.go`: Replaced two-step `MoveJump()` + `UpdateJump()` with single `MoveAndUpdateJump()` call; removed handler-level range validation (now in metadata)
- `server/handlers/jump_test.go`: Added `TestUpdateJump_MoveFailsDateValidation_Rollback` and `TestUpdateJump_MoveAndDateChange_Atomic` (26 handler tests total, all passing)

