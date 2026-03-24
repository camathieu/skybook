---
ticket: "004"
epic: jump-api
milestone: v1
title: Update Jump
status: done
priority: high
estimate: S
---

# Update Jump

`PUT /api/v1/jumps/:id` — Update an existing jump.

## Acceptance Criteria

- [x] Accepts JSON body with updatable fields (all except ID, UserID, CreatedAt)
- [x] Full replacement semantics (`PUT`) — all updatable fields must be provided
- [x] If `Number` differs from current, triggers a dedicated `MoveJump` DB transaction to reposition
- [x] `Number` must be between 1 and `MAX(Number)` — returns `400` if out of range
- [x] Validates required fields (Date, Dropzone) and JumpType enum
- [x] Returns `200 OK` with updated jump object
- [x] Returns `404 Not Found` / `400 Bad Request` as appropriate
- [x] Query scoped to current user (UserID=1 in v1)


## Done

- Epic completed: All 6 HTTP handlers implemented in `server/handlers/jump.go`
- Metadata layer extended with MoveJump and JumpFilters
- Fully tested (37 passing tests across handlers and metadata)
- Verified via `/review-changes`
