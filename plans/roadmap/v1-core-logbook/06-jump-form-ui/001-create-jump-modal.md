---
ticket: "001"
epic: jump-form-ui
milestone: v1
title: Create Jump Modal
status: done
priority: high
estimate: L
---

# Create Jump Modal

Modal form for creating a new jump, optimized for speed.

## Acceptance Criteria

- [x] Opens via `+ New Jump` button or `N` keyboard shortcut
- [x] Form fields matching Jump model — grouped logically (core, details, flags)
- [x] Required fields clearly marked: Date (defaults to today), Dropzone, JumpType
- [x] JumpType as dropdown with all enum values
- [x] Date picker with time component
- [x] Optional "Insert at position" to specify jump number
- [x] Submit creates jump via API, refreshes table, shows success toast
- [x] `Escape` or backdrop click closes with unsaved-changes warning
- [x] **Mobile responsive**: Modal becomes a full-screen sheet on screens < 640px.
- [x] **Touch-friendly**: Inputs have min 44px hit areas on mobile.
- [x] **Validation**: Errors displayed inline for missing required fields (Date, Dropzone, JumpType).

## Done

- `JumpModal.vue` created with Core / Details / People & Events / Flags / Notes sections
- `JumpList.vue` wired up `+ New Jump` button and global `N` keyboard shortcut with input guard
- Mobile full-screen sheet via `@media (max-width: 639px)` at the overlay level
- Date field defaults to `todayDatetime()` on open; inline validation on submit
- API errors displayed inline in `.save-error` block (no toast system yet — separate concern)
- Unsaved-changes guard via `isDirty` watcher + `confirm()` on close
