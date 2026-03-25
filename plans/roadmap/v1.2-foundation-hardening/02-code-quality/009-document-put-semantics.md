---
ticket: "009"
epic: code-quality
milestone: v1.2
title: Document PUT Semantics and Boolean Behavior
status: done
priority: medium
estimate: S
---

# Document PUT Semantics and Boolean Behavior

The API uses `PUT` for updates (full resource replacement), but Go's JSON unmarshaling makes boolean `false` indistinguishable from "field not sent" (both decode to `false`). This means if the frontend omits `nightJump`, it silently flips to `false`.

Currently the frontend always sends all fields, so this works. But it's undocumented and fragile — any future frontend change or third-party API consumer could trigger silent data loss.

## Acceptance Criteria

- [x] Document in `server/ARCHITECTURE.md` (§ PUT Semantics — UpdateJump):
  - PUT semantics require the client to send ALL fields
  - Boolean fields default to `false` if omitted (Go zero value behavior)
  - Immutable fields protected by `updatableColumns` whitelist
  - No PATCH support in v1.x

## Done

Already documented in `server/ARCHITECTURE.md` lines 146–155 (§ PUT Semantics — UpdateJump).
No additional code or documentation changes required. Closing as complete.
