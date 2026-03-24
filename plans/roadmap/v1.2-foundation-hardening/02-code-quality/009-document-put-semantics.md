---
ticket: "009"
epic: code-quality
milestone: v1.2
title: Document PUT Semantics and Boolean Behavior
status: planned
priority: medium
estimate: S
---

# Document PUT Semantics and Boolean Behavior

The API uses `PUT` for updates (full resource replacement), but Go's JSON unmarshaling makes boolean `false` indistinguishable from "field not sent" (both decode to `false`). This means if the frontend omits `nightJump`, it silently flips to `false`.

Currently the frontend always sends all fields, so this works. But it's undocumented and fragile — any future frontend change or third-party API consumer could trigger silent data loss.

## Acceptance Criteria

- [ ] Document in root `ARCHITECTURE.md` (§ API Contract):
  - PUT semantics require the client to send ALL fields
  - Boolean fields default to `false` if omitted (Go zero value behavior)
  - Optional numeric fields (`altitude`, `freefallTime`, `canopySize`) default to `null` if omitted (pointer type)
- [ ] Document in `webapp/ARCHITECTURE.md`:
  - `JumpModal.vue` `buildPayload()` MUST include all Jump fields
  - Explain why partial updates are dangerous with PUT
- [ ] Add a comment in `handlers/jump.go` `UpdateJump` explaining the PUT-requires-all-fields contract
- [ ] Consider: should the handler reject requests that are missing required boolean fields? (discuss in the ticket, decide at implementation time)
- [ ] Note in ARCHITECTURE.md that PATCH with field mask is a future option (deferred, not v1.2 scope)
