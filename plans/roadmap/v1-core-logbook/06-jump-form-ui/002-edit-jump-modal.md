---
ticket: "002"
epic: jump-form-ui
milestone: v1
title: Edit Jump Modal
status: done
priority: high
estimate: S
---

# Edit Jump Modal

Modal for editing an existing jump, reusing the create form component.

## Acceptance Criteria

- [x] Opens on table row click
- [x] Pre-populated with existing jump data
- [x] Same form layout and validation as Create
- [x] Number field editable (triggers reorder on save)
- [x] Shows jump number in header: "Edit Jump #142"
- [x] Submit updates via API, refreshes table
- [x] **Mobile responsive**: Modal becomes a full-screen sheet on screens < 640px.

## Done

- `JumpModal.vue` unified create/edit; `isEdit` computed from `jump` prop presence
- `JumpTable.vue` emits `edit` on row click/enter; `JumpCard.vue` emits `edit` on card click/enter
- `JumpList.vue` passes `editingJump` ref to `JumpModal :jump` prop
- `store.updateJump(id, payload)` calls `PUT /api/v1/jumps/:id` then re-fetches the list
