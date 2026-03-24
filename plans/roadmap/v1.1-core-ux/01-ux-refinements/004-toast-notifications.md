---
ticket: "004"
epic: ux-refinements
milestone: v1.1
title: Toast Notifications
status: planned
priority: medium
estimate: M
---

# Toast Notifications

A unified snackbar/toast system for non-blocking success and error feedback.

## Acceptance Criteria

- [ ] Global toast notification manager (e.g. Pinia store or simple composable).
- [ ] Toast component to display short, self-dismissing messages (e.g., 3-5 seconds).
- [ ] Success styling for actions like jump creation, editing, and deletion.
- [ ] Error styling to handle and display API errors gracefully in the UI.
- [ ] Stackable or replaceable toasts if multiple distinct events happen closely.
