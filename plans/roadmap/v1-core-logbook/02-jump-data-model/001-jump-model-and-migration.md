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

- [ ] `server/common/jump.go` — Jump struct with all v1 fields, GORM tags, JSON serialization
- [ ] All fields from PRD §3.1:

| Field | Go Type | GORM / Notes |
|-------|---------|-------------|
| `ID` | `uint` | PK, auto-increment |
| `UserID` | `uint` | FK, indexed — anonymous user (ID=1) in v1 |
| `Number` | `uint` | unique per user, indexed |
| `Date` | `time.Time` | required |
| `Dropzone` | `string` | required, indexed |
| `Aircraft` | `string` | optional |
| `JumpType` | `string` | required, indexed — enum validated at handler |
| `Altitude` | `*uint` | optional (pointer = nullable) defaults to 4000 |
| `FreefallTime` | `*uint` | optional, seconds |
| `CanopySize` | `*uint` | optional, sq ft |
| `Coach` | `string` | optional |
| `Event` | `string` | optional |
| `Description` | `string` | optional, text column |
| `Links` | `datatypes.JSONSlice[string]` | optional, stored as JSON text |
| `Landing` | `string` | optional — Stand-up / Sliding / PLF / Off-DZ / Water |
| `NightJump` | `bool` | default false |
| `OxygenJump` | `bool` | default false |
| `CutAway` | `bool` | default false |
| `CreatedAt` | `time.Time` | auto |
| `UpdatedAt` | `time.Time` | auto |

- [ ] `JumpType` enum constants defined in `jump.go`
- [ ] `UserID` foreign key (multi-tenant readiness)
- [ ] `Links` stored as JSON text, serialized/deserialized as `[]string`
- [ ] Migration via gormigrate creates `jumps` table
- [ ] Indexes: `(user_id, number)` unique, `date`, `dropzone`, `jump_type`
- [ ] Unit tests for JSON serialization of `Links`

## Technical Notes

- Use `gorm.io/datatypes` for `datatypes.JSONSlice[string]` (clean, no custom scanner needed)
- Nullable numeric fields use pointer types (`*uint`) so zero ≠ "not set"
- `Equipment`/Gear tracking deferred to v11 — full `Gear` table + `Kit` grouping system
- `Buddies` (many-to-many join table) deferred to v4; field omitted from this migration

