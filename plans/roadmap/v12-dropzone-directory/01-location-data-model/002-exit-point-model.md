---
ticket: "002"
epic: location-data-model
milestone: v12
title: ExitPoint Model
status: planned
priority: high
estimate: M
---

# ExitPoint Model

## Acceptance Criteria

- [ ] `ExitPoint` GORM model in `server/common/exit_point.go` per PRD §3.10
- [ ] Fields: Name, Object (`BUILDING`, `ANTENNA`, `SPAN`, `EARTH`, `OTHER`), Country, Region, Latitude, Longitude, Notes
- [ ] Globally shared table — **not** user-scoped
- [ ] Migration creates `exit_points` table
- [ ] Unit tests for model validation
