---
ticket: "003"
epic: jump-form-ui
milestone: v1
title: Autocomplete Fields
status: done
priority: medium
estimate: M
---

# Autocomplete Fields

Reusable autocomplete input component backed by the autocomplete API.

## Acceptance Criteria

- [x] Generic `AutocompleteInput.vue` component
- [x] Props: `field` (dropzone/aircraft/equipment/coach/event), `modelValue`, `placeholder`
- [x] Fetches suggestions from `GET /api/v1/jumps/autocomplete/:field?q=...`
- [x] Shows dropdown with top matches, ranked by usage count
- [x] Keyboard navigation: arrow keys, Enter to select, Escape to dismiss
- [x] Debounced input (200ms)
- [x] Allows freeform input (not restricted to existing values)
- [x] Dropdown floats above other modal content (z-index)
- [x] Closes dropdown on blur/outside click

## Done

- `AutocompleteInput.vue` created at `webapp/src/components/`
- 200ms debounce via `setTimeout`/`clearTimeout`; cleared in `onBeforeUnmount`
- Arrow keys traverse `activeIndex`; Enter selects; Escape dismisses
- `@blur` with 150ms delay lets `@mousedown.prevent` on suggestion items register first
- Suggestions list positioned `absolute` with `z-index: 1000` (inside JumpModal z-index: 1500)
- Touch targets fixed to `min-height: 44px` during review pass
