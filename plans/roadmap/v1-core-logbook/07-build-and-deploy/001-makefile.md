---
ticket: "001"
epic: build-and-deploy
milestone: v1
title: Makefile
status: planned
priority: high
estimate: M
---

# Makefile

Plik-style Makefile for building frontend, server, running tests, and dev mode.

## Acceptance Criteria

- [ ] `make frontend` — `cd webapp && npm ci && npm run build`
- [ ] `make server` — `go build` with embedded webapp, static linking, build info
- [ ] `make all` — frontend + server
- [ ] `make dev` — concurrent Vite dev server + Go server with live reload
- [ ] `make test` — Go tests
- [ ] `make test-frontend` — Vitest
- [ ] `make lint` — go fmt, go vet
- [ ] `make clean` — remove build artifacts
- [ ] Build tags: `osusergo,netgo,sqlite_omit_load_extension`
