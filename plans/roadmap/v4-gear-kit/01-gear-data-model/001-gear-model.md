---
ticket: "001"
epic: gear-data-model
milestone: v4
title: Gear Model
status: planned
priority: high
estimate: M
---

# Gear Model

## Acceptance Criteria

- [ ] `Gear` GORM model in `server/common/gear.go` per PRD §3.7
- [ ] Fields: Type, Manufacturer, Model, Size, Serial, DOM, PurchaseDate, PurchasePrice, NextMaintenanceAt, Notes, Active
- [ ] `GearType` enum with values: `CANOPY`, `RESERVE`, `HARNESS`, `AAD`, `HELMET`, `ALTIMETER`, `JUMPSUIT`, `WINGSUIT`, `CAMERA`, `AUDIBLE`, `OTHER`
- [ ] `GearType.IsValid()` validation method
- [ ] `AllGearTypes()` helper
- [ ] `UserID` FK for multi-tenant scoping
- [ ] Migration creates `gear` table
- [ ] Unit tests for GearType validation
