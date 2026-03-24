---
ticket: "002"
epic: jump-list-ui
milestone: v1
title: Search & Filters
status: done
priority: high
estimate: M
---

# Search & Filters

Filter bar above the jump table with text search and field-specific filters.

## Acceptance Criteria

- [x] Search input with `/` keyboard shortcut for focus
- [x] Filter dropdowns: jump type, dropzone, date range
- [x] Boolean toggles: cutaway, night
- [x] Filters update URL query params (shareable/bookmarkable)
- [x] Active filter chips with clear button
- [x] Debounced text search (300ms)

## Done

- Created `webapp/src/components/SearchBar.vue` — debounced 300ms search with `/` keyboard shortcut, clear button, focus ring
- Created `webapp/src/components/FilterBar.vue` — jump type dropdown, dropzone autocomplete dropdown (fetched from `/api/v1/jumps/autocomplete/dropzone`), date range inputs, boolean toggles (cutaway/night with 3-state: null/true/false), active filter chips with individual remove and "Clear all"
- Filters sync to URL query params via store `toQuery()`/`initFromQuery()` for bookmarkable/shareable URLs
