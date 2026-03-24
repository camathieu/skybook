---
ticket: "006"
epic: ux-refinements
milestone: v1.1
title: Save & Add Another
status: planned
priority: high
estimate: S
---

# Save & Add Another

Allow users to quickly log sequential jumps without leaving the jump modal.

## Acceptance Criteria

- [ ] Add a "Save & Add Another" button next to "Save" in `JumpModal.vue`.
- [ ] Clicking it triggers the standard jump creation via the API.
- [ ] Upon success, instead of closing the modal, it clears only jump-specific fields (e.g., Altitude, Freefall Time, Description) while retaining session-specific defaults (e.g., Date, Dropzone, Aircraft, JumpType).
- [ ] Refreshes the underlying jump table/list in the background so the new jump is visible.
