---
ticket: "004"
epic: backend-foundation
milestone: v1
title: Server Skeleton
status: planned
priority: high
estimate: S
---

# Server Skeleton

Set up the HTTP server with gorilla/mux, middleware chain, and health endpoint.

## Acceptance Criteria

- [ ] `server/server/` package with router setup
- [ ] Middleware chain: recovery, logging, request ID
- [ ] `GET /health` returns `200 OK`
- [ ] `GET /api/v1/config` returns server config (feature flags)
- [ ] Graceful shutdown on SIGTERM/SIGINT
- [ ] Listen address/port from config

## Technical Notes

- Plik-style middleware chain composition
- JSON error responses with consistent format: `{"error": "message"}`
