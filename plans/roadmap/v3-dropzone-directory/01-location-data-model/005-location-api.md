---
ticket: "005"
epic: location-data-model
milestone: v3
title: Location API
status: planned
priority: high
estimate: L
---

# Location API

## Acceptance Criteria

- [ ] `server/metadata/dropzone.go` — CRUD + autocomplete for Dropzone
- [ ] `server/metadata/exit_point.go` — CRUD + autocomplete for ExitPoint
- [ ] `server/metadata/wind_tunnel.go` — CRUD + autocomplete for WindTunnel
- [ ] `server/handlers/location.go` — HTTP handlers for all three types

### Dropzone endpoints
- [ ] `GET /api/v1/dropzones` — list / search dropzones
- [ ] `GET /api/v1/dropzones/autocomplete` — search by name (prefix match)
- [ ] `POST /api/v1/dropzones` — create dropzone
- [ ] `GET /api/v1/dropzones/:id` — get single
- [ ] `PUT /api/v1/dropzones/:id` — update
- [ ] `DELETE /api/v1/dropzones/:id` — delete

### ExitPoint endpoints
- [ ] `GET /api/v1/exit-points` — list / search
- [ ] `GET /api/v1/exit-points/autocomplete` — search by name
- [ ] `POST /api/v1/exit-points` — create
- [ ] `GET /api/v1/exit-points/:id` — get single
- [ ] `PUT /api/v1/exit-points/:id` — update
- [ ] `DELETE /api/v1/exit-points/:id` — delete

### WindTunnel endpoints
- [ ] `GET /api/v1/wind-tunnels` — list / search
- [ ] `GET /api/v1/wind-tunnels/autocomplete` — search by name
- [ ] `POST /api/v1/wind-tunnels` — create
- [ ] `GET /api/v1/wind-tunnels/:id` — get single
- [ ] `PUT /api/v1/wind-tunnels/:id` — update
- [ ] `DELETE /api/v1/wind-tunnels/:id` — delete

### Registration
- [ ] All routes registered in `server/server/server.go`
- [ ] Unit tests for all handlers
