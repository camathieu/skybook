---
ticket: "002"
epic: ux-refinements
milestone: v1.1
title: Custom Filter Dropdowns (Alphabetical)
status: done
priority: medium
estimate: S
---

# Custom Filter Dropdowns (Alphabetical)

Replace the native browser `<select>` dropdowns in the filter bar with a custom-styled Vue component (similar to Plik's theme selector), sorted alphabetically, limited to ~10 visible items with a scrollbar. Add Aircraft as a new filter option alongside Dropzone and Jump Type.

## Acceptance Criteria

- [x] **Backend**: Add `Aircraft` to `metadata.JumpFilters` and apply it in `GetJumps()`.
- [x] **Backend**: Read `aircraft` query parameter in `handlers.ListJumps()`.
- [x] **Backend API**: Add optional `?sort=alpha` to `handlers.Autocomplete()` / `metadata.GetJumpAutocomplete()` to order by `col ASC`. (Modal still defaults to recency).
- [x] **Frontend Component**: Create `CustomSelect.vue` that renders a styled, absolutely positioned dropdown menu with a dark theme design, a checkmark for the selected item, max-height (scrollable after ~10 items), and proper click-outside handling.
- [x] **Frontend Integration**: Update `FilterBar.vue` to replace the native `<select>` dropdowns for Dropzone, Jump Type, and Aircraft with `CustomSelect.vue`.
- [x] **Frontend State**: Add `aircraft` to `stores/jumps.js` filters and pass it in `fetchJumps()`.
- [x] **Data Fetching**: FilterBar calls API for Dropzone, Aircraft, and JumpType with `?sort=alpha` so the lists are perfectly alphabetical.

## Done

- `webapp/src/components/CustomSelect.vue` (new): Dark-themed custom dropdown with animation, checkmark, max-height scroll, click-outside close, and full mobile support
- `server/metadata/jump.go`: Added `Aircraft` to `JumpFilters`; applied filter in `GetJumps()`; `GetJumpAutocomplete` now takes `sortBy` string (`"alpha"` → `col ASC`); `jump_type` added to `allowedAutocompleteFields`
- `server/handlers/jump.go`: `ListJumps` reads `aircraft` param; `Autocomplete` reads `sort` param
- `server/handlers/jump_test.go`: Added `TestAutocomplete_AlphaSort` test
- `webapp/src/stores/jumps.js`: Added `aircraft` to `filters`, `resetFilters`, `hasActiveFilters`, and `fetchJumps` params
- `webapp/src/components/FilterBar.vue`: Replaced native `<select>` elements with `CustomSelect` for Dropzone, Aircraft (new), and JumpType; hardcoded `jumpTypes` array removed
