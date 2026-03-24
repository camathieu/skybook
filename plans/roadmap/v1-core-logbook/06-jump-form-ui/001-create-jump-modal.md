---
ticket: "001"
epic: jump-form-ui
milestone: v1
title: Create Jump Modal
status: planned
priority: high
estimate: L
---

# Create Jump Modal

Modal form for creating a new jump, optimized for speed.

## Acceptance Criteria

- [ ] Opens via `+ New Jump` button or `N` keyboard shortcut
- [ ] Form fields matching Jump model — grouped logically (core, details, flags)
- [ ] Required fields clearly marked: Date (defaults to today), Dropzone, JumpType
- [ ] JumpType as dropdown with all enum values
- [ ] Date picker with time component
- [ ] Optional "Insert at position" to specify jump number
- [ ] Submit creates jump via API, refreshes table, shows success toast
- [ ] `Escape` or backdrop click closes with unsaved-changes warning
