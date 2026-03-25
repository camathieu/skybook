---
ticket: "008"
epic: code-quality
milestone: v1.2
title: Complete or Remove Aircraft Filter
status: done
priority: low
estimate: S
---

# Complete or Remove Aircraft Filter

The backend `JumpFilters` struct has an `Aircraft` field and the metadata layer applies it in queries, but:
- The frontend `FilterBar.vue` doesn't expose an aircraft filter dropdown
- The Pinia store's `initFromQuery`/`toQuery` doesn't handle an `aircraft` query param

This is a dangling wire — the backend supports it but it's unreachable from the UI.

## Acceptance Criteria

**Option A — Complete it (chosen):**
- [x] Add an aircraft filter dropdown to `FilterBar.vue` (using the existing autocomplete endpoint)
- [x] Wire `aircraft` through `initFromQuery` and `toQuery` in the jumps store

## Done

The aircraft filter was already fully implemented (Option A) during the Jump List UI epic. The
only genuine gap was URL persistence:

- `webapp/src/stores/jumps.js` — `initFromQuery`: restores `aircraft` from URL query param
- `webapp/src/stores/jumps.js` — `toQuery`: persists `aircraft` to URL query param

All other wiring (`FilterBar.vue` dropdown, `fetchJumps`, `resetFilters`, `hasActiveFilters`) was
already in place.
