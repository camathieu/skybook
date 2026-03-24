---
ticket: "004"
epic: data-integrity
milestone: v1.2
title: Selective Column Update
status: planned
priority: high
estimate: M
---

# Selective Column Update

`UpdateJump` uses GORM `.Save(jump)` which overwrites **every column** with whatever the client sent. This means:
- Zero-valued fields silently clear existing data (e.g. `nightJump: false` is indistinguishable from "field not sent")
- As the model grows (buddies v4, documents v5), any new field not sent by the frontend will be wiped

This is a time bomb that gets worse with every field added to the Jump model.

## Acceptance Criteria

- [ ] Replace `tx.Save(jump)` in `UpdateJump` with `tx.Select(fields...).Save(jump)` or `tx.Model(jump).Updates(map)` using an explicit whitelist of updatable columns
- [ ] The whitelist includes only user-mutable fields: `date`, `dropzone`, `aircraft`, `jump_type`, `altitude`, `freefall_time`, `canopy_size`, `lo`, `event`, `description`, `links`, `landing`, `night_jump`, `oxygen_jump`, `cut_away`, `packjob`
- [ ] Immutable fields are never overwritten: `id`, `user_id`, `number`, `created_at`
- [ ] `updated_at` is still auto-managed by GORM
- [ ] Test: update a jump, verify that `created_at` is unchanged
- [ ] Test: send a partial body (missing `nightJump`), verify existing `nightJump=true` is preserved

## Design Note

The partial-body test (#6) requires either:
- **Option A**: PATCH semantics with field mask (complex, but correct)
- **Option B**: Frontend always sends all fields (current behavior, enforce via test)

For v1.2, Option B is acceptable. Document this clearly in ARCHITECTURE.md. PATCH can be added in a later milestone if needed.
