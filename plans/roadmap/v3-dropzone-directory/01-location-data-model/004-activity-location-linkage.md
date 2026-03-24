---
ticket: "004"
epic: location-data-model
milestone: v3
title: Activity Location Linkage
status: planned
priority: high
estimate: M
---

# Activity Location Linkage

## Acceptance Criteria

- [ ] Add `DropzoneID *uint` nullable FK to `Jump` → `Dropzone`
- [ ] Add `ExitPointID *uint` nullable FK to `BaseJump` → `ExitPoint`
- [ ] Add `WindTunnelID *uint` nullable FK to `TunnelSession` → `WindTunnel`
- [ ] Existing freeform string fields (`Dropzone`, `Location`, `Tunnel`) kept as fallback — not removed
- [ ] Migration adds FK columns to existing tables
- [ ] Jump/BaseJump/TunnelSession CRUD handlers preload location entity on read
- [ ] Create/update handlers accept optional location ID and sync the FK
- [ ] Unit tests for FK linkage on all three activity types
