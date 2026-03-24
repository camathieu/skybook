---
ticket: "005"
epic: jump-api
milestone: v1
title: Delete Jump (Renumber)
status: done
priority: high
estimate: S
---

# Delete Jump (Renumber)

`DELETE /api/v1/jumps/:id` — Delete a jump and renumber subsequent jumps downward.

## Acceptance Criteria

- [x] Deletes the jump by ID
- [x] All jumps with `Number > deleted.Number` shift down by 1
- [x] Wrapped in a transaction
- [x] Returns `204 No Content` on success
- [x] Returns `404 Not Found` if jump doesn't exist
- [x] Query scoped to current user (UserID=1 in v1)


## Done

- Epic completed: All 6 HTTP handlers implemented in `server/handlers/jump.go`
- Metadata layer extended with MoveJump and JumpFilters
- Fully tested (37 passing tests across handlers and metadata)
- Verified via `/review-changes`
