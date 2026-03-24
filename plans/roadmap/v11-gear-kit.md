---
milestone: v11
title: Gear & Kit Tracking
status: planned
---

# v11 — Gear & Kit Tracking

Full equipment tracking system for skydivers. Track individual gear items (canopy, reserve, harness, AAD, etc.) and group them into kits for one-click assignment to jumps.

## Epics

- [Gear Data Model](v11-gear-kit/01-gear-data-model.md) — Gear/Kit models, jump linkage, API
- [Gear UI](v11-gear-kit/02-gear-ui.md) — Equipment management page, kit management, jump form integration

## Tickets

### Gear Data Model
- [001 — Gear Model](v11-gear-kit/01-gear-data-model/001-gear-model.md)
- [002 — Kit Model](v11-gear-kit/01-gear-data-model/002-kit-model.md)
- [003 — Jump Gear Linkage](v11-gear-kit/01-gear-data-model/003-jump-gear-linkage.md)
- [004 — Gear API](v11-gear-kit/01-gear-data-model/004-gear-api.md)

### Gear UI
- [001 — Gear Management Page](v11-gear-kit/02-gear-ui/001-gear-management-page.md)
- [002 — Kit Management](v11-gear-kit/02-gear-ui/002-kit-management.md)
- [003 — Jump Form Gear Selector](v11-gear-kit/02-gear-ui/003-jump-form-gear-selector.md)

## Overview

### Gear

User-scoped equipment items with type, manufacturer, model, size, serial, DOM, purchase info, and maintenance tracking.

**GearType enum**: `CANOPY` · `RESERVE` · `HARNESS` · `AAD` · `HELMET` · `ALTIMETER` · `JUMPSUIT` · `WINGSUIT` · `CAMERA` · `AUDIBLE` · `OTHER`

### Kit

A convenience grouping of gear items (e.g. "Main rig" = Canopy + Harness + Reserve + AAD). Selecting a kit on the jump form pre-fills gear items — **no FK on the Jump table**, kits are purely a UI shortcut.

### Jump linkage

`Jump` → `Gear` via `jump_gear` many-to-many join table. Individual gear items can be added/removed freely.

## Links

- PRD §3.7 (Gear), §3.8 (Kit)
- Supersedes `Wingsuit bool` on `Jump`
