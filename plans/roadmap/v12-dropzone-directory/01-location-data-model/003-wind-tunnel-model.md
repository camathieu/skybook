---
ticket: "003"
epic: location-data-model
milestone: v12
title: WindTunnel Model
status: planned
priority: high
estimate: S
---

# WindTunnel Model

## Acceptance Criteria

- [ ] `WindTunnel` GORM model in `server/common/wind_tunnel.go` per PRD §3.11
- [ ] Fields: Name, Country, City, DiameterFt, Website, Email, Phone
- [ ] Globally shared table — **not** user-scoped
- [ ] Migration creates `wind_tunnels` table
- [ ] Unit tests for model validation
