---
ticket: "002"
epic: backend-foundation
milestone: v1
title: Config System
status: planned
priority: high
estimate: S
---

# Config System

Implement TOML-based configuration with environment variable overrides using `SKYBOOK_` prefix (Plik pattern).

## Acceptance Criteria

- [ ] `server/common/config.go` — Config struct with TOML tags
- [ ] Default `skybook.cfg` template file with commented options
- [ ] `SKYBOOK_` prefixed env var overrides (screaming snake case)
- [ ] Config sections: `[server]`, `[database]`, `[defaults]`
- [ ] Unit tests for config loading and env override

## Technical Notes

- Use `BurntSushi/toml` for parsing (same as Plik)
- Auth and storage config sections can be stubbed for v1
