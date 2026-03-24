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

- [ ] Extract a `validateJumpFields(jump *common.Jump) (string, bool)` helper (returns error message + ok)
- [ ] Or move validation into the metadata layer (e.g. `jump.Validate() error`)
- [ ] Both `CreateJump` and `UpdateJump` handlers use the shared validation
- [ ] Consider adding validation for fields that currently have none:
  - `altitude` upper bound (e.g. ≤ 50,000 ft — world record is ~41,000 ft)
  - `freefallTime` upper bound (e.g. ≤ 600 seconds — 10 min max)
  - `canopySize` upper bound (e.g. ≤ 500 sq ft)
  - `landing` enum validation if a fixed set is expected
- [ ] Existing handler tests pass
- [ ] New tests for any added validation bounds
