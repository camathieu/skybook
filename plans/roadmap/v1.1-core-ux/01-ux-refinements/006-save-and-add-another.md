---
ticket: "006"
epic: ux-refinements
milestone: v1.1
title: Save & Add Another
status: done
priority: high
estimate: S
---

# Save & Add Another

Allow users to quickly log sequential jumps without leaving the jump modal.

## Acceptance Criteria

- [x] Add a "Save & Add Another" button next to "Save" in `JumpModal.vue`.
- [x] Clicking it triggers the standard jump creation via the API.
- [x] Upon success, instead of closing the modal, it clears only jump-specific fields (e.g., Altitude, Freefall Time, Description) while retaining session-specific defaults (e.g., Date, Dropzone, Aircraft, JumpType).
- [x] Refreshes the underlying jump table/list in the background so the new jump is visible.

## Done

- Added `saveAndAddAnother()` function in `JumpModal.vue`
- Shows success toast with jump number via `useToastStore`
- Clears: number, altitude, freefallTime, canopySize, lo, event, landing, flags, description
- Retains: date, dropzone, aircraft, jumpType
- Re-focuses date field after clear
- Button only visible in create mode (`v-if="!isEdit"`)
