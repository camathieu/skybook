---
ticket: "013"
epic: code-quality
milestone: v1.2
title: Server-Driven Frontend Settings
status: planned
priority: medium
estimate: M
---

# Server-Driven Frontend Settings

## Context

Frontend constants (altitude presets, landing options, pattern presets, jump types) are currently hardcoded in `JumpModal.vue` and bundled at build time by Vite. This means they can't be changed without rebuilding the SPA.

Following the Plik pattern (`settings.json` served by the backend), create a server-side settings endpoint that serves these constants at runtime. The Go server reads them from `skybook.cfg` (with sensible defaults) and exposes them via an API endpoint. The Vue app fetches them on init and uses them instead of hardcoded arrays.

This enables runtime configuration without rebuilding the frontend — the single Go binary ships with defaults, and operators can override presets in `skybook.cfg`.

## Acceptance Criteria

### Backend
- [ ] Add a `[presets]` section to `skybook.cfg` with default values for:
  - `jump_types` — discipline list (default: current JUMP_TYPES)
  - `landing_options` — landing presets (default: current LANDING_OPTIONS)
  - `pattern_options` — canopy pattern presets (default: current PATTERN_PRESETS)
  - `altitude_presets_ft` — feet altitude presets (default: current ALTITUDE_PRESETS)
  - `altitude_presets_m` — meters altitude presets (default: `['1500', '3000', '3600', '4000', '4500', '6000']`)
- [ ] Parse the `[presets]` section in `common/config.go`
- [ ] Serve `GET /api/v1/settings` returning the presets as JSON (no auth required)
- [ ] Environment variable overrides work (e.g. `SKYBOOK_PRESETS_LANDING_OPTIONS="Stand-up,Sliding,PLF"`)

### Frontend
- [ ] Fetch `/api/v1/settings` on app init (in `App.vue` or a Pinia store)
- [ ] Replace all hardcoded constant arrays in `JumpModal.vue` with store-driven values
- [ ] Graceful fallback to hardcoded defaults if the endpoint is unreachable

### Documentation
- [ ] Document the `[presets]` config section in `skybook.cfg` with comments
- [ ] Update `ARCHITECTURE.md` docs
