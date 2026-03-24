---
ticket: "003"
epic: jump-list-ui
milestone: v1
title: Pagination
status: done
priority: medium
estimate: S
---

# Pagination

Page controls for navigating through the jump list.

## Acceptance Criteria

- [x] Page indicator: "Showing 1–25 of 342 jumps"
- [x] Previous / Next buttons with disabled state at boundaries
- [x] Per-page selector: 25 / 50 / 100
- [x] URL query param sync (`?page=2&per_page=50`)

## Done

- Created `webapp/src/components/Pagination.vue` — "Showing X–Y of Z jumps" indicator, Prev/Next buttons with disabled boundaries, per-page selector (25/50/100) as segmented button group
- Page and per_page sync to URL query params; changing filters auto-resets to page 1
- Touch-friendly mobile sizing (44px min-height on all buttons)
