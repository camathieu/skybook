---
ticket: "003"
epic: jump-api
milestone: v1
title: Get Jump
status: done
priority: high
estimate: XS
---

# Get Jump

`GET /api/v1/jumps/:id` — Get a single jump by ID.

## Acceptance Criteria

- [x] Returns `200 OK` with full jump object
- [x] Returns `404 Not Found` if jump doesn't exist
- [x] Query scoped to current user (UserID=1 in v1)


## Done

- Epic completed: All 6 HTTP handlers implemented in `server/handlers/jump.go`
- Metadata layer extended with MoveJump and JumpFilters
- Fully tested (37 passing tests across handlers and metadata)
- Verified via `/review-changes`
