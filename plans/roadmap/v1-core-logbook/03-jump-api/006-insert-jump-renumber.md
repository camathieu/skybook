---
ticket: "006"
epic: jump-api
milestone: v1
title: Insert Jump (Renumber)
status: done
priority: high
estimate: S
---

# Insert Jump (Renumber)

`POST /api/v1/jumps` with `number` field — Insert a jump at a specific position.

## Acceptance Criteria

- [x] When `number` is provided in the request body, insert at that position
- [x] All jumps with `Number >= requested` shift up by 1
- [x] Wrapped in a transaction
- [x] Validates that requested number is between 1 and `MAX(Number) + 1`
- [x] Returns `201 Created` with the jump at the requested position
- [x] UserID hardcoded to anonymous user (ID=1) in v1


## Done

- Epic completed: All 6 HTTP handlers implemented in `server/handlers/jump.go`
- Metadata layer extended with MoveJump and JumpFilters
- Fully tested (37 passing tests across handlers and metadata)
- Verified via `/review-changes`
