---
ticket: "003"
epic: testing
milestone: v1
title: E2E Tests
status: done
priority: medium
estimate: M
---

# E2E Tests

Playwright E2E tests for critical user flows.

## Acceptance Criteria

- [x] Create a jump, verify it appears in the table
- [x] Edit a jump, verify changes persist
- [x] Delete a jump, verify renumbering
- [x] Insert a jump at position, verify numbering
- [x] Search and filter jumps
- [x] Responsive layout on mobile viewport

## Done
- Installed `@playwright/test` and Chromium browser binary
- Created `playwright.config.js` with Desktop Chrome and Mobile Safari (iPhone SE) projects
- `webServer` config auto-boots Go backend (`:memory:` SQLite) + Vite dev server before running suite
- Added `data-testid` attributes to all key elements in `JumpList.vue`, `JumpModal.vue`, `ConfirmModal.vue`, `SearchBar.vue`, `JumpTable.vue`, `JumpCard.vue`
- Created `e2e/jumps.spec.js` — full CRUD flow tests (create, edit, delete, insert-at, search)
- Created `e2e/mobile.spec.js` — responsive layout, nav, touch-target, and modal usability tests at 375px
- Runnable via `make test-e2e`
