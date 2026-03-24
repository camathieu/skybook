---
milestone: v3
title: Location Directory
status: planned
---

# v3 — Location Directory

Replace freeform location strings on all activity types with references to shared canonical tables: **Dropzone** (jumps), **ExitPoint** (BASE), **WindTunnel** (tunnel sessions).

## Epics

- [Location Data Model](v3-dropzone-directory/01-location-data-model.md) — Dropzone, ExitPoint, WindTunnel models, linkage, API, seed data
- [Location UI](v3-dropzone-directory/02-location-ui.md) — Autocomplete, management pages, location resolver

## Tickets

### Location Data Model
- [001 — Dropzone Model](v3-dropzone-directory/01-location-data-model/001-dropzone-model.md)
- [002 — ExitPoint Model](v3-dropzone-directory/01-location-data-model/002-exit-point-model.md)
- [003 — WindTunnel Model](v3-dropzone-directory/01-location-data-model/003-wind-tunnel-model.md)
- [004 — Activity Location Linkage](v3-dropzone-directory/01-location-data-model/004-activity-location-linkage.md)
- [005 — Location API](v3-dropzone-directory/01-location-data-model/005-location-api.md)
- [006 — Seed Data](v3-dropzone-directory/01-location-data-model/006-seed-data.md)

### Location UI
- [001 — Location Autocomplete](v3-dropzone-directory/02-location-ui/001-location-autocomplete.md)
- [002 — Location Management Pages](v3-dropzone-directory/02-location-ui/002-location-management-pages.md)
- [003 — Activity Location Resolver](v3-dropzone-directory/02-location-ui/003-activity-location-resolver.md)

## Overview

### Dropzone
Globally shared table (not user-scoped). Fields: Name, Country, City, ICAO, GPS, Website, Email, Phone. Linked to `Jump` via nullable `DropzoneID` FK.

### ExitPoint
Globally shared. Fields: Name, Object type (B.A.S.E.), Country, Region, GPS, Notes. Linked to `BaseJump` via nullable `ExitPointID` FK.

### WindTunnel
Globally shared. Fields: Name, Country, City, DiameterFt, Website, Email, Phone. Linked to `TunnelSession` via nullable `WindTunnelID` FK.

### Migration strategy

1. Add `dropzones`, `exit_points`, `wind_tunnels` tables + seed community data
2. Add nullable FK columns to `jumps`, `base_jumps`, `tunnel_sessions`
3. Existing freeform strings kept as fallback
4. UI "Link" action to resolve string → FK
5. Eventually deprecate freeform strings (future milestone)

## Links

- PRD §3.9 (Dropzone), §3.10 (ExitPoint), §3.11 (WindTunnel)
- Supersedes `Dropzone string` on `Jump`, `Location string` on `BaseJump`, `Tunnel string` on `TunnelSession`
