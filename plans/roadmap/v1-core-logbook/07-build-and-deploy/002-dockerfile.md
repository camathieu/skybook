---
ticket: "002"
epic: build-and-deploy
milestone: v1
title: Dockerfile
status: done
priority: medium
estimate: S
---

# Dockerfile

Multi-stage Dockerfile for building and running SkyBook.

## Acceptance Criteria

- [x] Stage 1: Node — build frontend (`npm ci && npm run build`)
- [x] Stage 2: Go — build server with embedded frontend
- [x] Stage 3: Alpine — minimal runtime image with the single binary
- [x] Exposes port 8080, volume for `skybook.db`
- [x] `.dockerignore` for node_modules, .git, build artifacts

## Done
- Wrote multi-stage `Dockerfile` with node and golang builder
- Deployed statically linked binary to Alpine runtime container
- Checked `.dockerignore` against common artifacts
