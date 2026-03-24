---
ticket: "002"
epic: jump-api
milestone: v1
title: List Jumps (Paginated)
status: done
priority: high
estimate: M
---

# List Jumps (Paginated)

`GET /api/v1/jumps` — List jumps with pagination, sorting, and filtering.

## Acceptance Criteria

- [x] Pagination: `page` (1-based), `per_page` (default 25, max 100)
- [x] `per_page` values above 100 are silently clamped to 100
- [x] Sorting: `sort` (number, date, dropzone, altitude) + `order` (asc/desc)
- [x] Default sort: `number desc` when no sort/order params provided
- [x] Filters: `q` (full-text), `date_from`, `date_to`, `dropzone`, `jump_type`, `altitude_min`, `altitude_max`, `cutaway`, `night`, `lo`
- [x] Response includes: `items`, `total`, `page`, `per_page`, `total_pages`
- [x] Returns `200 OK` with empty items array when no matches
- [x] Full-text search `q` searches across description, dropzone, event, lo (case-insensitive)
- [x] All queries scoped to current user (UserID=1 in v1)

## Technical Notes

- GORM scopes pattern for composable filters
- `q` filter uses SQLite `LIKE` with `%query%` (case-insensitive via `COLLATE NOCASE`)


## Done

- Epic completed: All 6 HTTP handlers implemented in `server/handlers/jump.go`
- Metadata layer extended with MoveJump and JumpFilters
- Fully tested (37 passing tests across handlers and metadata)
- Verified via `/review-changes`
