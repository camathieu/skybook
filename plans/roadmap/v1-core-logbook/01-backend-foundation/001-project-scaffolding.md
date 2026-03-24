---
ticket: "001"
epic: backend-foundation
milestone: v1
title: Project Scaffolding
status: done
priority: high
estimate: S
---

# Project Scaffolding

Initialize the Go module and directory structure following the Plik-style layered architecture.

## Acceptance Criteria

- [x] `go mod init github.com/root-gg/skybook`
- [x] Directory tree: `server/{common,metadata,handlers,middleware,server,cmd}`
- [x] `server/main.go` entrypoint with cobra root command
- [x] `.gitignore` for Go, Node, SQLite artifacts
- [x] Empty `webapp/` directory with `.gitkeep`
- [x] `server/cmd/root.go` with cobra setup, `--config` flag
- [x] `server/cmd/serve.go` as the default command (starts the server)

## Done

- Created `server/go.mod` — module `github.com/root-gg/skybook`, Go 1.26
- Created `server/main.go` — 3-line entry delegating to cobra
- Created `server/cmd/root.go` — cobra root with `--config` flag (default `skybook.cfg`)
- Created `server/cmd/serve.go` — default serve command: loads config → validate → init DB → start server → graceful shutdown on SIGTERM/SIGINT
- Created `.gitignore` — covers Go binaries, Node modules, SQLite files, IDE configs
- Created `webapp/.gitkeep` — placeholder for webapp directory
