---
ticket: "007"
epic: ux-refinements
milestone: v1.1
title: Packjob Boolean
status: done
priority: medium
estimate: M
---

# Packjob Boolean

Track whether a parachute was packed by a packer for a given jump, serving as a reminder to pay them.

## Acceptance Criteria

- [x] Add a `Packjob` boolean field to the Go `Jump` model (maps to DB).
- [x] Database migration handled via GORM `AutoMigrate` (pre-v1 policy).
- [x] Ensure the API allows creating and updating jumps with the `Packjob` field.
- [x] Add a "Packjob" toggle or checkbox to the Flags section of `JumpModal.vue`.
- [ ] (Optional) Add a visual indicator in the `JumpList` or `JumpCard` to show if a jump was a packjob.

## Done

- Added `Packjob bool` to `server/common/jump.go` with `gorm:"default:false" json:"packjob"`
- Added `packjob` to `form` state, `buildPayload()`, and `saveAndAddAnother()` clear list in `JumpModal.vue`
- Added 📦 Packjob toggle in the Flags section of `JumpModal.vue`
- Lint, Go build, frontend build, 55 unit tests all pass
