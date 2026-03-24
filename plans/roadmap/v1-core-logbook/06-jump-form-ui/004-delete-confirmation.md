---
ticket: "004"
epic: jump-form-ui
milestone: v1
title: Delete Confirmation
status: in-progress
priority: medium
estimate: XS
---

# Delete Confirmation

Confirmation dialog before deleting a jump.

## Acceptance Criteria

- [ ] Delete button in edit modal footer
- [ ] Confirmation dialog: "Delete Jump #142? This will renumber all subsequent jumps."
- [ ] Danger-styled confirm button
- [ ] Calls DELETE API, refreshes table, shows success toast
- [ ] Dialog is styled as a custom UI modal
- [ ] Loading state on the confirm button while deletion is in progress
