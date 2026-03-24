---
ticket: "004"
epic: webapp-foundation
milestone: v1
title: SPA Embedding
status: planned
priority: high
estimate: S
---

# SPA Embedding

Embed the built webapp into the Go binary using `embed.FS` and serve it as a SPA.

## Acceptance Criteria

- [ ] `//go:embed all:dist` directive in the server package (or a dedicated `webapp/embed.go`)
- [ ] SPA handler: serves static files from `dist/`, falls back to `index.html` for client-side routes
- [ ] API routes (`/api/*`) take precedence over SPA catch-all
- [ ] Cache headers: hashed assets get long cache, `index.html` gets no-cache
- [ ] Works with `make server` after `make frontend`
