---
ticket: "004"
epic: jump-form-ui
milestone: v1
title: Delete Confirmation
status: done
priority: medium
estimate: XS
---

# Delete Confirmation

Confirmation dialog before deleting a jump.

## Acceptance Criteria

- [x] Delete button in edit modal footer
- [x] Confirmation dialog: "Delete Jump #142? This will renumber all subsequent jumps."
- [x] Danger-styled confirm button
- [x] Calls DELETE API, refreshes table, shows success toast
- [x] Dialog is styled as a custom UI modal
- [x] Loading state on the confirm button while deletion is in progress

## Done

- `ConfirmModal.vue` created as a generic reusable danger dialog (Teleport to body)
- Used inside `JumpModal.vue` via `v-if="showDeleteConfirm"` with `title`, `message`, `confirm-text`, `:danger="true"`, `:loading="deleting"` props
- `store.deleteJump(id)` calls `DELETE /api/v1/jumps/:id` then re-fetches list (handles renumbering)
- `.spinner` CSS animation shown on confirm button while `deleting` is true
- Dialog has `z-index: 2000` (above JumpModal's 1500)
