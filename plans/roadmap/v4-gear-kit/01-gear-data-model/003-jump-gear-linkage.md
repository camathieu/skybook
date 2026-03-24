---
ticket: "003"
epic: gear-data-model
milestone: v4
title: Jump Gear Linkage
status: planned
priority: high
estimate: M
---

# Jump Gear Linkage

## Acceptance Criteria

- [ ] Add `GearItems []Gear` many-to-many on `Jump` via `jump_gear` join table
- [ ] Migration creates `jump_gear` table
- [ ] Migration removes `Wingsuit bool` from `jumps` table (superseded by `GearType=WINGSUIT`)
- [ ] Jump CRUD handlers preload `GearItems` on read
- [ ] Jump create/update handlers accept `gearItemIds []uint` and sync the association
- [ ] Unit tests for gear association on Jump create/update/delete
