---
milestone: v12
title: Dropzone Directory
status: planned
---

# v12 — Dropzone Directory

Replace the `Dropzone string` field on `Jump` (and `TunnelSession`) with a reference to a shared `Dropzone` table, enabling richer location data and cross-user disambiguation.

## Epics

- [01 — Dropzone Data Model](v12-dropzone-directory/01-dropzone-data-model.md)
- [02 — Dropzone UI](v12-dropzone-directory/02-dropzone-ui.md)

## Overview

### Dropzone Table

A globally shared table (not scoped per user — one canonical entry per real-world dropzone):

| Field | Type | Description |
|-------|------|-------------|
| `ID` | `uint` (PK) | |
| `Name` | `string` | Canonical name, e.g. "Skydive Empuriabrava" |
| `Country` | `string` | ISO country code |
| `City` | `string` | Nearest city |
| `ICAO` | `string` | Optional ICAO airport code |
| `Latitude` | `float64` | Optional GPS |
| `Longitude` | `float64` | Optional GPS |

### Jump linkage

`Jump.DropzoneID` nullable FK → `Dropzone`. During migration, existing `Dropzone string` values are matched against the new table (fuzzy match + user confirmation in UI). The raw string is kept as a fallback until confirmed.

### Autocomplete upgrade

Dropzone autocomplete queries the shared table instead of `DISTINCT(dropzone)` on the jumps table — allowing search across all known dropzones, not just visited ones.

## Migration strategy

1. Add `dropzones` table
2. Add `DropzoneID *uint` nullable FK to `jumps`
3. Import seed data (optional community-maintained list)
4. UI: "Link" button on each jump to resolve the string → FK
5. Eventually make `DropzoneID` required (v12+)

## Links

- PRD §3.x (to be added)
- Supersedes `Dropzone string` on `Jump`, `BaseJump`, `TunnelSession`
