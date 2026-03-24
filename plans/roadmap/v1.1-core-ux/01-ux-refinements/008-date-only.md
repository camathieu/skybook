---
ticket: "008"
epic: ux-refinements
milestone: v1.1
title: Date-Only Field (Drop Time)
status: planned
priority: high
estimate: S
---

# Date-Only Field (Drop Time)

Change the jump date from a `datetime` to a `date`-only field. Skydivers care about the day they jumped, not the exact time.

## Acceptance Criteria

- [ ] Backend: Change `Date` field in the `Jump` model from `time.Time` with time precision to date-only storage (`DATE` type in SQLite, or `time.Time` truncated to midnight UTC)
- [ ] Backend: API accepts dates in `YYYY-MM-DD` format (e.g., `"2025-03-08"`)
- [ ] Backend: API returns dates in `YYYY-MM-DD` format (no time component)
- [ ] Frontend: `JumpModal.vue` uses `<input type="date">` instead of `<input type="datetime-local">`
- [ ] Frontend: Jump list/cards display dates as `Mar 8, 2025` (no time)
- [ ] Migration: Existing datetime values are truncated to date (no data loss risk since time is not meaningful)
- [ ] `fakedb` generates date-only values
- [ ] All tests updated to use date-only format
