---
ticket: "016"
epic: code-quality
milestone: v1.2
title: Replace Browser Alert with Custom Unsaved Changes Modal
status: planned
priority: medium
estimate: S
---

# Replace Browser Alert with Custom Unsaved Changes Modal

## Context

The JumpModal currently uses `window.confirm("You have unsaved changes. Close anyway?")` when the user tries to close with unsaved edits. This native browser dialog looks ugly and inconsistent with the app's dark theme. Replace it with a custom styled confirmation modal that matches the existing design system.

## Acceptance Criteria

- [ ] Replace `window.confirm()` with a custom confirmation modal component
- [ ] Modal matches the app's dark theme and design system
- [ ] Two buttons: "Discard" (closes without saving) and "Keep editing" (returns to form)
- [ ] Modal is accessible: focus trap, Escape key support, aria attributes
- [ ] Touch targets ≥ 44px for mobile
