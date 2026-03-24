---
ticket: "003"
epic: jump-form-ui
milestone: v1
title: Autocomplete Fields
status: planned
priority: medium
estimate: M
---

# Autocomplete Fields

Reusable autocomplete input component backed by the autocomplete API.

## Acceptance Criteria

- [ ] Generic `AutocompleteInput.vue` component
- [ ] Props: `field` (dropzone/aircraft/equipment/coach/event), `modelValue`, `placeholder`
- [ ] Fetches suggestions from `GET /api/v1/jumps/autocomplete/:field?q=...`
- [ ] Shows dropdown with top matches, ranked by usage count
- [ ] Keyboard navigation: arrow keys, Enter to select, Escape to dismiss
- [ ] Debounced input (200ms)
- [ ] Allows freeform input (not restricted to existing values)
