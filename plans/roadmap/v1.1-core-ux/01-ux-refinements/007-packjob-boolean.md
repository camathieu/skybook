---
ticket: "007"
epic: ux-refinements
milestone: v1.1
title: Packjob Boolean
status: planned
priority: medium
estimate: M
---

# Packjob Boolean

Track whether a parachute was packed by a packer for a given jump, serving as a reminder to pay them.

## Acceptance Criteria

- [ ] Add a `Packjob` boolean field to the Go `Jump` model (maps to DB).
- [ ] Add a database migration for the new field.
- [ ] Ensure the API allows creating and updating jumps with the `Packjob` field.
- [ ] Add a "Packjob" toggle or checkbox to the Flags section of `JumpModal.vue`.
- [ ] (Optional) Add a visual indicator in the `JumpList` or `JumpCard` to show if a jump was a packjob.
