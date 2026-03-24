---
ticket: "003"
epic: backend-foundation
milestone: v1
title: Database Setup
status: planned
priority: high
estimate: S
---

# Database Setup

Initialize GORM with SQLite in WAL mode, set up migration framework.

## Acceptance Criteria

- [ ] `server/metadata/` package with GORM setup
- [ ] SQLite connection with WAL mode and busy timeout
- [ ] Gormigrate-based migration runner
- [ ] Database path from config (`[database].Path`)
- [ ] Auto-create database file and parent directories
- [ ] Unit tests for database initialization

## Technical Notes

- Use `gorm.io/driver/sqlite` with CGO (same build tags as Plik: `sqlite_omit_load_extension`)
- WAL mode: `PRAGMA journal_mode=WAL; PRAGMA busy_timeout=5000;`
