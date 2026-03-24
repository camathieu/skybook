---
ticket: "004"
epic: testing
milestone: v1
title: FakeDB Subcommand
status: planned
priority: medium
estimate: M
---

# FakeDB Subcommand

A `skybook fakedb` cobra subcommand that generates a pre-populated SQLite database with randomised test data. Useful for UI development, manual testing, and performance benchmarking — same pattern as `plikd fakedb`.

## Acceptance Criteria

- [ ] `skybook fakedb` subcommand in `server/cmd/fakedb.go`
- [ ] Flags: `--users`, `--jumps-per-user`, `--output` (default `/tmp/test-skybook.db`)
- [ ] Creates an admin user with known credentials (login: `admin`, password: `skybook`)
- [ ] Generates randomised jumps with realistic field values (dates, altitudes, freefall times, disciplines, locations, aircraft, equipment)
- [ ] Jump numbers follow the contiguous 1-based sequence invariant per user
- [ ] Progress logging during generation (every N users)
- [ ] Summary output at completion (user count, jump count, elapsed time)
- [ ] Prints usage instructions showing how to start the server with the generated database
- [ ] Works with the existing GORM backend and migration system
