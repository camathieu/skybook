---
ticket: "008"
epic: code-quality
milestone: v1.2
title: Complete or Remove Aircraft Filter
status: planned
priority: low
estimate: S
---

# Complete or Remove Aircraft Filter

The backend `JumpFilters` struct has an `Aircraft` field and the metadata layer applies it in queries, but:
- The frontend `FilterBar.vue` doesn't expose an aircraft filter dropdown
- The Pinia store's `initFromQuery`/`toQuery` doesn't handle an `aircraft` query param

This is a dangling wire — the backend supports it but it's unreachable from the UI.

## Acceptance Criteria

**Option A — Complete it:**
- [ ] Add an aircraft filter dropdown to `FilterBar.vue` (using the existing autocomplete endpoint)
- [ ] Wire `aircraft` through `initFromQuery` and `toQuery` in the jumps store
- [ ] Add frontend test for the aircraft filter

**Option B — Remove it:**
- [ ] Remove `Aircraft` from `JumpFilters`
- [ ] Remove the `aircraft` filter clause from `GetJumps` in metadata
- [ ] Remove the `aircraft` query param parsing from the handler

Choose whichever option makes sense at implementation time. If other filter UX improvements are planned, Option A is preferred.
