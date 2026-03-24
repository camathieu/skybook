---
ticket: "007"
epic: jump-api
milestone: v1
title: Autocomplete Endpoint
status: done
priority: medium
estimate: S
---

# Autocomplete Endpoint

`GET /api/v1/jumps/autocomplete/:field` — Return distinct values for a field, ranked by frequency.

## Acceptance Criteria

- [x] Supported fields: `dropzone`, `aircraft`, `equipment`, `lo`, `event`
- [x] Returns top 20 values sorted by usage count descending
- [x] Accepts `q` query param for prefix filtering (case-insensitive)
- [x] Returns `400 Bad Request` for unsupported field names
- [x] Response: `[{"value": "Skydive DeLand", "count": 42}, ...]`
- [x] Query scoped to current user (UserID=1 in v1)


## Done

- Epic completed: All 6 HTTP handlers implemented in `server/handlers/jump.go`
- Metadata layer extended with MoveJump and JumpFilters
- Fully tested (37 passing tests across handlers and metadata)
- Verified via `/review-changes`
