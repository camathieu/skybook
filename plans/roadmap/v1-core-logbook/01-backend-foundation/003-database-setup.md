---
ticket: "003"
epic: backend-foundation
milestone: v1
title: Database Setup
status: done
priority: high
estimate: S
---

# Database Setup

Initialize GORM with SQLite in WAL mode, set up migration framework.

## Acceptance Criteria

- [x] `server/metadata/` package with GORM setup
- [x] SQLite connection with WAL mode and busy timeout
- [x] Gormigrate-based migration runner
- [x] Database path from config (`[database].Path`)
- [x] Auto-create database file and parent directories
- [x] Unit tests for database initialization
- [x] `Shutdown()` method to close the DB connection cleanly

## Done

- `server/metadata/backend.go` — `Backend` struct with `NewBackend()`, `DB()`, `Shutdown()`, `migrate()`. SQLite DSN includes `_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=ON`. Auto-creates parent directories with `os.MkdirAll`. Skips migration when list is empty (gormigrate requires at least one entry).
- `server/metadata/migrations.go` — empty migration list, ready for v1 Jump model
- `server/metadata/backend_test.go` — 4 tests: init, WAL mode verification, auto-dir creation, shutdown
