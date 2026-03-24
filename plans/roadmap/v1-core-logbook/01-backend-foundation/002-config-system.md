---
ticket: "002"
epic: backend-foundation
milestone: v1
title: Config System
status: done
priority: high
estimate: S
---

# Config System

Implement TOML-based configuration with environment variable overrides using `SKYBOOK_` prefix.

## Acceptance Criteria

- [x] `server/common/config.go` — Config struct with TOML tags
- [x] Default `server/skybook.cfg` template file with commented options
- [x] `SKYBOOK_` prefixed env var overrides (screaming snake case)
- [x] Config sections: `[server]`, `[database]`, `[defaults]`
- [x] Unit tests for config loading and env override
- [x] Config validation (port range, non-empty database path, valid unit system)
- [x] `InitializeDefaults()` method to fill zero values with sensible defaults

## Done

- `server/common/config.go` — `Config`, `ServerConfig`, `DatabaseConfig`, `DefaultsConfig` structs with TOML tags; `LoadConfig()`, `NewConfig()`, `InitializeDefaults()`, `Validate()`, `ApplyEnvironment()` using reflection
- `server/common/http.go` — `WriteJSON()`/`WriteError()` JSON response helpers
- `server/skybook.cfg` — commented TOML template covering all options plus future stubs
- `server/common/config_test.go` — 5 tests: defaults, TOML load, env overrides, validation, missing file
