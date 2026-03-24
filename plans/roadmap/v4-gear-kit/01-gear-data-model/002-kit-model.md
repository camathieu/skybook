---
ticket: "002"
epic: gear-data-model
milestone: v4
title: Kit Model
status: planned
priority: high
estimate: S
---

# Kit Model

## Acceptance Criteria

- [ ] `Kit` GORM model in `server/common/gear.go` (or `kit.go`) per PRD §3.8
- [ ] Fields: Name, GearItems (many-to-many via `kit_gear` join table)
- [ ] `UserID` FK for multi-tenant scoping
- [ ] Migration creates `kits` and `kit_gear` tables
- [ ] No FK from Kit to Jump — kits are a UI convenience only
- [ ] Unit tests for Kit ↔ Gear association
