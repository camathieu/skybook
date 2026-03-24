---
ticket: "002"
epic: build-and-deploy
milestone: v1
title: Dockerfile
status: planned
priority: medium
estimate: S
---

# Dockerfile

Multi-stage Dockerfile for building and running SkyBook.

## Acceptance Criteria

- [ ] Stage 1: Node — build frontend (`npm ci && npm run build`)
- [ ] Stage 2: Go — build server with embedded frontend
- [ ] Stage 3: Alpine — minimal runtime image with the single binary
- [ ] Exposes port 8080, volume for `skybook.db`
- [ ] `.dockerignore` for node_modules, .git, build artifacts
