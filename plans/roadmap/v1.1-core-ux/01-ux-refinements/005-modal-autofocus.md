---
ticket: "005"
epic: ux-refinements
milestone: v1.1
title: Modal Autofocus
status: done
priority: medium
estimate: XS
---

# Modal Autofocus

Automatically focus the primary input field when opening the jump modal to speed up data entry.

## Acceptance Criteria

- [x] When `JumpModal.vue` is opened in "Create" mode, automatically focus the `#f-date` Date & Time field via `onMounted` + `nextTick`.
- [x] Guarded with `window.matchMedia('(min-width: 640px)')` — autofocus is skipped on mobile to avoid virtual keyboard pop-up.

## Done

- Added `dateRef` template ref to the `#f-date` input in `JumpModal.vue`
- Added `onMounted` with `nextTick` to focus date field on create mode (desktop only)
- Build and 55 unit tests pass
- Browser-verified: `document.activeElement.id === 'f-date'` immediately on modal open
