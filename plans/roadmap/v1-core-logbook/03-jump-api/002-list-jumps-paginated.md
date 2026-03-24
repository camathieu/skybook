---
ticket: "002"
epic: jump-api
milestone: v1
title: List Jumps (Paginated)
status: planned
priority: high
estimate: M
---

# List Jumps (Paginated)

`GET /api/v1/jumps` — List jumps with pagination, sorting, and filtering.

## Acceptance Criteria

- [ ] Pagination: `page` (1-based), `per_page` (default 25, max 100)
- [ ] Sorting: `sort` (number, date, dropzone, altitude) + `order` (asc/desc, default: desc by number)
- [ ] Filters: `q` (full-text), `date_from`, `date_to`, `dropzone`, `jump_type`, `altitude_min`, `altitude_max`, `cutaway`, `night`, `coach`
- [ ] Response includes: `items`, `total`, `page`, `per_page`, `total_pages`
- [ ] Returns `200 OK` with empty items array when no matches
- [ ] Full-text search `q` searches across description, dropzone, event, coach

## Technical Notes

- GORM scopes pattern for composable filters
- `q` filter uses SQLite `LIKE` with `%query%` (case-insensitive via `COLLATE NOCASE`)
