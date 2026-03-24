---
ticket: "004"
epic: webapp-foundation
milestone: v1
title: SPA Embedding
status: done
priority: high
estimate: S
---

# SPA Embedding

Embed the built webapp into the Go binary using `embed.FS` and serve it as a SPA.

## Acceptance Criteria

- [x] `//go:embed all:dist` directive in the server package (or a dedicated `webapp/embed.go`)
- [x] SPA handler: serves static files from `dist/`, falls back to `index.html` for client-side routes
- [x] API routes (`/api/*`) take precedence over SPA catch-all
- [x] Cache headers: hashed assets get long cache, `index.html` gets no-cache
- [x] Works with `make server` after `make frontend`

## Done

- Added `spa.go` to `server/server` package, utilizing `//go:embed all:dist`.
- Handled static file serving with aggressive caching for hashed assets, while falling back to `index.html` properly for client routing.
- Added `Makefile` with targets for building and dev servers, plus copying the `dist` folder to satisfy the `go:embed` on `make server`.
- Created `server/server/dist/.gitkeep` and added negation rule to `.gitignore` to enable bare Go builds without failing the embed command.
