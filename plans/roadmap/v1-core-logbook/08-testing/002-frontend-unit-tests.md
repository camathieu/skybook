---
ticket: "002"
epic: testing
milestone: v1
title: Frontend Unit Tests
status: done
priority: medium
estimate: S
---

# Frontend Unit Tests

Vitest unit tests for Vue components and stores.

## Acceptance Criteria

- [x] Jump store tests (CRUD actions, state management)
- [x] AutocompleteInput component tests
- [x] Jump form validation tests (via store)
- [x] Pagination component tests

## Done
- Installed `vitest`, `@vue/test-utils`, `jsdom`, `@vitest/coverage-v8`
- Configured Vitest in `vite.config.js` (jsdom env, scoped to `src/**/*.spec.js`)
- Created `src/stores/jumps.spec.js` — 14 tests covering state init, CRUD, sorting, filtering, URL sync, and pagination
- Created `src/components/AutocompleteInput.spec.js` — 6 tests covering rendering, v-model, debounce, keyboard events
- Created `src/components/Pagination.spec.js` — 10 tests covering page range display, button states, per-page controls
- All 30 tests pass via `make test-frontend`
