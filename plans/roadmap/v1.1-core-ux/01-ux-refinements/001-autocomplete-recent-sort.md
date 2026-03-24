---
ticket: "001"
epic: ux-refinements
milestone: v1.1
title: Autocomplete Recent Sort
status: done
priority: medium
estimate: S
---

# Autocomplete Recent Sort

Refine the autocomplete API and component to sort suggestions primarily by recency of use, and show suggestions on field focus (not only when typing).

## Acceptance Criteria

- [x] Backend: Change `/api/v1/jumps/autocomplete/:field` sort from `COUNT(*) DESC` to `MAX(date) DESC` (most recently used first)
- [x] Backend: When no `q` query parameter is provided, return all distinct non-empty values (sorted by recency) — this powers the on-focus dropdown
- [x] Frontend: `AutocompleteInput.vue` shows suggestions on focus even when the field is empty (loads recent values immediately)
- [x] Frontend: Selecting a suggestion from the dropdown on focus works identically to the typing flow
- [x] Frontend: Clearing a field and re-focusing also re-shows all suggestions (intentional UX — the dropdown acts as a "recent values" list)
- [x] API response is simplified to a plain `[]string` — removes the `{value, count}` object shape
- [x] Existing autocomplete unit tests pass and new test cases cover: recency sort order, on-focus suggestions

## Done

- `server/metadata/jump.go`: Removed `AutocompleteResult` struct; `GetJumpAutocomplete` now returns `[]string` sorted by `MAX(date) DESC` + `col ASC` secondary
- `server/handlers/jump.go`: `Autocomplete` handler updated to write plain `[]string` (+ new `sort` param for ticket 002)
- `server/handlers/jump_test.go`: `TestAutocomplete_Dropzone` updated to check `[]string` + recency ordering (Perris first)
- `webapp/src/components/AutocompleteInput.vue`: Added `onFocus` handler; shared `fetchSuggestions()`; documented clear-field re-fetch behavior
- `webapp/src/components/FilterBar.vue`: Removed `.map(r => r.value)` mapping (API now returns plain strings)
- `webapp/src/components/AutocompleteInput.spec.js`: Updated tests + added on-focus test
