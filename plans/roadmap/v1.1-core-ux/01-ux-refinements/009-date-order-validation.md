---
ticket: "009"
epic: ux-refinements
milestone: v1.1
title: Date Order Validation
status: planned
priority: high
estimate: S
depends_on: ["008"]
---

# Date Order Validation

Enforce that jump dates are monotonically non-decreasing by jump number. Multiple jumps on the same day are allowed, but jump #N can never have a date *after* jump #N+1, or *before* jump #N-1.

## Acceptance Criteria

- [ ] Backend: On `CreateJump` — validate date ≥ previous jump's date (if inserting, also ≤ next jump's date)
- [ ] Backend: On `UpdateJump` — if date is changed, validate against neighbors (jump #N-1 and #N+1)
- [ ] Backend: Return `400 Bad Request` with a clear error message: `"date 2025-03-01 is before jump #7 (2025-03-05)"`
- [ ] Frontend: Display the validation error in the modal (inline or toast)
- [ ] Same-day jumps are always valid (comparison is `>=` / `<=`, not strict)
- [ ] When appending a new jump (no insert-at), date must be ≥ the last jump's date
- [ ] Edge cases: first jump (no predecessor) and single-jump logbook always pass validation
- [ ] Tests cover: append valid, append invalid (date before last), insert-at valid, insert-at invalid (both neighbors), update date valid, update date invalid
