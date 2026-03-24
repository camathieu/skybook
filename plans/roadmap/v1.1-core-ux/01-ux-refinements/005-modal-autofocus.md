---
ticket: "005"
epic: ux-refinements
milestone: v1.1
title: Modal Autofocus
status: planned
priority: medium
estimate: XS
---

# Modal Autofocus

Automatically focus the primary input field when opening the jump modal to speed up data entry.

## Acceptance Criteria

- [ ] When `JumpModal.vue` is opened in "Create" mode (e.g. via `+ New Jump` or `N` shortcut), automatically focus the "Dropzone" or first active form field.
- [ ] Test that this autofocus behavior does not trigger virtual keyboards aggressively on mobile devices unless intended.
