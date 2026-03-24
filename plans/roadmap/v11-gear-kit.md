---
milestone: v11
title: Gear & Kit Tracking
status: planned
---

# v11 — Gear & Kit Tracking

Full equipment tracking system for skydivers. Track individual gear items and group them into kits for convenience.

## Epics

- [01 — Gear Data Model](v11-gear-kit/01-gear-data-model.md)
- [02 — Gear UI](v11-gear-kit/02-gear-ui.md)

## Overview

### Gear Items

A `Gear` table shared across all user-scoped activity types (jumps, BASE, tunnel as applicable). Each item has:

| Field | Type | Description |
|-------|------|-------------|
| `ID` | `uint` (PK) | |
| `UserID` | `uint` (FK) | |
| `Type` | `string` | `CANOPY`, `HARNESS`, `RESERVE`, `AAD`, `WINGSUIT`, `HELMET`, `CAMERA`, `SUIT`, `OTHER` |
| `Name` | `string` | e.g. "Pilot 168", "Skyhook", "Vigil 2" |
| `Notes` | `string` | Optional freeform notes |
| `Active` | `bool` | Whether currently in use |

### Kits

A `Kit` groups gear items for convenience — e.g. "My main rig" = Canopy + Harness + Reserve + AAD.

| Field | Type | Description |
|-------|------|-------------|
| `ID` | `uint` (PK) | |
| `UserID` | `uint` (FK) | |
| `Name` | `string` | e.g. "Race rig", "Student rig" |
| `GearItems` | `[]Gear` | Many-to-many via `kit_gear` join table |

### Jump linkage

Jumps get a `KitID` nullable FK and/or a `GearItems []Gear` many-to-many. Selecting a kit pre-populates gear for a jump; individual items can be overridden.

## Links

- PRD §3.x (to be added)
- Supersedes the `Wingsuit bool` field on `Jump` (migration to add `GearItems` and remove `Wingsuit`)
