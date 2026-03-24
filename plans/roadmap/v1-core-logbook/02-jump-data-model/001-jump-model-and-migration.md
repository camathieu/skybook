---
ticket: "001"
epic: jump-data-model
milestone: v1
title: Jump Model & Migration
status: planned
priority: high
estimate: M
---

# Jump Model & Migration

Define the Jump GORM model with all v1 fields and create the initial migration.

## Acceptance Criteria

- [ ] `server/common/jump.go` — Jump struct with GORM tags and JSON serialization
- [ ] All fields from PRD §3.1: Number, Date, Dropzone, Aircraft, JumpType, Altitude, DeployAltitude, FreefallTime, CanopySize, Coach, Event, Description, Links, Landing, NightJump, OxygenJump, CutAway
- [ ] `UserID` foreign key (for multi-tenant readiness)
- [ ] `Links` stored as JSON text, serialized/deserialized as `[]string`
- [ ] Migration via gormigrate creates `jumps` table
- [ ] Indexes on: `Number` (unique per user), `Date`, `Dropzone`, `JumpType`

## Technical Notes

- `JumpType` is a string enum validated at the handler level, not DB level
- `Links` uses GORM's `datatypes.JSON` or a custom scanner/valuer
