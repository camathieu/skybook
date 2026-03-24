---
ticket: "001"
epic: build-and-deploy
milestone: v1
title: Makefile
status: done
priority: high
estimate: M
---

# Makefile

Plik-style Makefile for building frontend, server, running tests, and dev mode.

## Acceptance Criteria

- [x] `make frontend` — `cd webapp && npm ci && npm run build`
- [x] `make server` — `go build` with embedded webapp, static linking, build info
- [x] `make all` — frontend + server
- [x] `make dev` — concurrent Vite dev server + Go server with live reload
- [x] `make test` — Go tests
- [x] `make test-frontend` — Vitest
- [x] `make lint` — go fmt, go vet
- [x] `make clean` — remove build artifacts
- [x] Build tags: `osusergo,netgo,sqlite_omit_load_extension`

## Done
- Updated Makefile with CGO tags for sqlite_omit_load_extension
- Set -ldflags for static builds
- Confirmed targets all run perfectly
