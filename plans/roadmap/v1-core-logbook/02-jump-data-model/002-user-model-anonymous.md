---
ticket: "002"
epic: jump-data-model
milestone: v1
title: User Model (Anonymous)
status: planned
priority: high
estimate: S
---

# User Model (Anonymous)

Create the User model and auto-provision an anonymous default user for v1 single-user mode.

## Acceptance Criteria

- [ ] `server/common/user.go` — User struct with all PRD §3.6 fields
- [ ] Migration creates `users` table
- [ ] On first startup, auto-create anonymous user (ID=1, Provider="local", Name="Skydiver")
- [ ] All jumps in v1 are attributed to this anonymous user
- [ ] `UnitSystem` defaults to "imperial"

## Technical Notes

- This establishes multi-tenant readiness: when v6 adds auth, no schema changes needed
- The anonymous user check runs in the migration or server startup sequence
