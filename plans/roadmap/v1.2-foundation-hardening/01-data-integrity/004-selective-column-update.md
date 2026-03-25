---
ticket: "004"
epic: data-integrity
milestone: v1.2
title: Selective Column Update
status: done
priority: high
estimate: M
---

# Selective Column Update

`UpdateJump` uses GORM `.Save(jump)` which overwrites **every column** with whatever the client sent. This means:
- Zero-valued fields silently clear existing data (e.g. `nightJump: false` is indistinguishable from "field not sent")
- As the model grows (buddies v4, documents v5), any new field not sent by the frontend will be wiped

This is a time bomb that gets worse with every field added to the Jump model.

## Acceptance Criteria

- [x] Replace `tx.Save(jump)` in `UpdateJump` with `tx.Select(fields...).Save(jump)` or `tx.Model(jump).Updates(map)` using an explicit whitelist of updatable columns
- [x] The whitelist includes only user-mutable fields: `date`, `dropzone`, `aircraft`, `jump_type`, `altitude`, `freefall_time`, `canopy_size`, `lo`, `event`, `description`, `links`, `landing`, `night_jump`, `oxygen_jump`, `cut_away`, `packjob`
- [x] Immutable fields are never overwritten: `id`, `user_id`, `number`, `created_at`
- [x] `updated_at` is still auto-managed by GORM
- [x] Test: update a jump, verify that `created_at` is unchanged
- [x] Test: send a partial body (missing `nightJump`), verify existing `nightJump=true` is preserved

## Design Note

`PUT /api/v1/jumps/:id` is full-replacement: the frontend always sends all user-mutable fields. The switch from `.Save()` to `.Select(whitelist).Updates()` is **not** about adding partial-update semantics — it is purely defensive: it hard-codes which columns can ever be modified, so `id`, `user_id`, `number`, and `created_at` cannot be accidentally overwritten regardless of what arrives in the struct.

PATCH with a field mask would be needed only if the frontend needed to send a subset of fields. That is not the case today and is deferred to a later milestone if required.

## Done

- Added `updatableColumns` whitelist (16 user-mutable fields + `updated_at`) in `server/metadata/jump.go`
- Replaced `tx.Omit(clause.Associations).Save(jump)` with `tx.Model(jump).Select(updatableColumns).Updates(jump)` in `updateJumpTx`; removed unused `clause` import
- Added `TestUpdateJump_CreatedAtUnchanged` and `TestUpdateJump_BooleanOverwrite` in `server/metadata/backend_test.go`
- Documented PUT full-replacement semantics in `server/ARCHITECTURE.md`
- All 6 test packages pass; lint and build clean
