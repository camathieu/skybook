---
ticket: "004"
epic: backend-foundation
milestone: v1
title: Server Skeleton
status: done
priority: high
estimate: S
---

# Server Skeleton

Set up the HTTP server with gorilla/mux, middleware chain, and health endpoint.

## Acceptance Criteria

- [x] `server/server/` package with router setup
- [x] Middleware chain: recovery, logging, request ID
- [x] `GET /health` returns `200 OK`
- [x] `GET /api/v1/config` returns server config (feature flags)
- [x] Graceful shutdown on SIGTERM/SIGINT
- [x] Listen address/port from config
- [x] `common.WriteError(w, message, code)` JSON error helper
- [x] `common.WriteJSON(w, data, code)` JSON response helper

## Done

- `server/server/server.go` — `SkyBookServer` with `NewSkyBookServer()`, `setupRouter()`, `Start()` (net.Listen + Serve), `Shutdown()` (10s context timeout)
- `server/middleware/recovery.go` — panic recovery → 500 JSON with stack trace logged
- `server/middleware/logging.go` — request logging (method, path, status, duration) using `statusWriter` to capture status code
- `server/middleware/request_id.go` — UUID request ID in `X-Request-ID` header + context
- `server/handlers/misc.go` — `HealthHandler` (→ `{"status":"ok"}`), `ConfigHandler` (→ public-safe defaults subset)
- Middleware order: Recovery → RequestID → Logging (recovery outermost for maximum safety)
