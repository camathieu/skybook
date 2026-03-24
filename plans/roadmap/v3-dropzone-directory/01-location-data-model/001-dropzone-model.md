---
ticket: "001"
epic: location-data-model
milestone: v3
title: Dropzone Model
status: planned
priority: high
estimate: M
---

# Dropzone Model

## Acceptance Criteria

- [ ] `Dropzone` GORM model in `server/common/dropzone.go` per PRD §3.9
- [ ] Fields: Name, Country (ISO 3166-1 alpha-2), City, ICAO, Latitude, Longitude, Website, Email, Phone
- [ ] Globally shared table — **not** user-scoped (no `UserID` FK)
- [ ] Migration creates `dropzones` table
- [ ] Unit tests for model validation
