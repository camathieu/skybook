---
ticket: "006"
epic: code-quality
milestone: v1.2
title: Shared Validation Helper
status: done
priority: medium
estimate: S
---

# Shared Validation Helper

`CreateJump` and `UpdateJump` handlers both validate the same 3 fields (date, dropzone, jumpType) with identical code blocks. This should be a single shared function.

## Acceptance Criteria

- [x] Extract a `validateJumpFields(jump *common.Jump) error` helper or add `Validate() error` to the `Jump` model
- [x] Both `CreateJump` and `UpdateJump` handlers use the shared validation
- [x] Existing handler tests pass
- [x] New tests for the validation logic (e.g., `TestJump_Validate`)

## Done

- Added `Jump.Validate() error` in `server/common/jump.go` (date, dropzone, jumpType checks)
- Replaced both 12-line duplicate validation blocks in `server/handlers/jump.go` with `jump.Validate()` calls; removed unused `strings` import
- Added `TestJump_Validate` in `server/common/jump_test.go` with 4 table-driven cases (valid, missing date, empty dropzone, invalid type)
- All 6 test packages pass

