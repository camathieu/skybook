---
ticket: "001"
epic: jump-api
milestone: v1
title: Create Jump
status: done
priority: high
estimate: M
---

# Create Jump

`POST /api/v1/jumps` — Create a new jump, appended at the end by default.

## Acceptance Criteria

- [x] Accepts JSON body with all Jump fields (except ID, Number, CreatedAt, UpdatedAt)
- [x] Auto-assigns `Number = MAX(Number) + 1`
- [x] Validates required fields: Date, Dropzone, JumpType
- [x] Validates JumpType against allowed enum values
- [x] Returns `201 Created` with the full jump object (including assigned Number)
- [x] Returns `400 Bad Request` with descriptive error for invalid input
- [x] UserID hardcoded to anonymous user (ID=1) in v1

## Technical Notes

- Number assignment uses the invariant logic from ticket `jump-data-model/003`
- UserID is hardcoded to anonymous user (ID=1) in v1


## Done

- Epic completed: All 6 HTTP handlers implemented in `server/handlers/jump.go`
- Metadata layer extended with MoveJump and JumpFilters
- Fully tested (37 passing tests across handlers and metadata)
- Verified via `/review-changes`
