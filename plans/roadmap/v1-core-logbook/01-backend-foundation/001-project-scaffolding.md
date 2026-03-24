---
ticket: "001"
epic: backend-foundation
milestone: v1
title: Project Scaffolding
status: planned
priority: high
estimate: S
---

# Project Scaffolding

Initialize the Go module and directory structure following the Plik-style layered architecture.

## Acceptance Criteria

- [ ] `go mod init github.com/root-gg/skybook`
- [ ] Directory tree: `server/{common,metadata,handlers,middleware,server,cmd}`
- [ ] `server/main.go` entrypoint with cobra root command
- [ ] `.gitignore` for Go, Node, SQLite artifacts
- [ ] Empty `webapp/` directory with `.gitkeep`

## Technical Notes

- Mirror the Plik package layering: `common → metadata → middleware → handlers → cmd → server`
- Use cobra for CLI commands from the start (even if v1 only has `serve`)
