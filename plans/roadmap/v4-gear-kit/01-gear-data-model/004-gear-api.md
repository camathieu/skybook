---
ticket: "004"
epic: gear-data-model
milestone: v4
title: Gear API
status: planned
priority: high
estimate: M
---

# Gear API

## Acceptance Criteria

- [ ] `server/metadata/gear.go` — CRUD operations for Gear and Kit
- [ ] `server/handlers/gear.go` — HTTP handlers

### Gear endpoints
- [ ] `GET /api/v1/gear` — list user's gear (filterable by type, active status)
- [ ] `POST /api/v1/gear` — create gear item
- [ ] `GET /api/v1/gear/:id` — get single gear item
- [ ] `PUT /api/v1/gear/:id` — update gear item
- [ ] `DELETE /api/v1/gear/:id` — delete gear item

### Kit endpoints
- [ ] `GET /api/v1/kits` — list user's kits (with preloaded gear items)
- [ ] `POST /api/v1/kits` — create kit (with gear item IDs)
- [ ] `GET /api/v1/kits/:id` — get single kit with gear items
- [ ] `PUT /api/v1/kits/:id` — update kit (name + gear item IDs)
- [ ] `DELETE /api/v1/kits/:id` — delete kit

### Registration
- [ ] Routes registered in `server/server/server.go`
- [ ] Unit tests for all handlers
