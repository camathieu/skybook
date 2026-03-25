---
ticket: "006"
epic: code-quality
milestone: v1.2
title: Shared Validation Helper
status: planned
priority: medium
estimate: S
---

# Shared Validation Helper

`CreateJump` and `UpdateJump` handlers both validate the same 3 fields (date, dropzone, jumpType) with identical code blocks. This should be a single shared function.

## Acceptance Criteria

- [ ] Extract a `validateJumpFields(jump *common.Jump) error` helper or add `Validate() error` to the `Jump` model
- [ ] Both `CreateJump` and `UpdateJump` handlers use the shared validation
- [ ] Existing handler tests pass
- [ ] New tests for the validation logic (e.g., `TestJump_Validate`)

