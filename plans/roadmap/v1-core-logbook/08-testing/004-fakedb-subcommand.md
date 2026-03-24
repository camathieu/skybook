---
ticket: "004"
epic: testing
milestone: v1
title: FakeDB Subcommand
status: done
priority: medium
estimate: M
---

# FakeDB Subcommand

A `skybook fakedb` cobra subcommand that generates a pre-populated SQLite database with randomised test data. Useful for UI development, manual testing, and performance benchmarking — same pattern as `plikd fakedb`.

## Acceptance Criteria

- [ ] `skybook fakedb` subcommand in `server/cmd/fakedb.go`
- [ ] Flags: `--jumps` (default `2000`), `--output` (default `skybook.db`)
- [ ] Generates jumps chronologically over the last 10 years (~200 jumps/year)
- [ ] Groups jumps by active days (3 to 8 jumps per day) at a consistent Dropzone and Aircraft per day
- [ ] Simulates 80% local DZ days and 20% weekend event days at random destination DZs
- [ ] Discipline distribution: 75% FF (Freefly), 20% WS (Wingsuit), 1% Hop & Pop, 4% random
- [ ] Logical correlations: Hop & Pops are 4k-5k altitude with 0-5s freefall; WS are 60-90s freefall; FF are 40-50s freefall
- [ ] Approximately 25% of jumps are marked as `Packjob = true`
- [ ] Progression: Equipment string changes over the 10 years (e.g., student rigs down to cross-braced canopies)
- [ ] Progression: Sprinkles 1 or 2 `Cutaway = true` jumps randomly across the logbook
- [ ] Jump numbers follow the contiguous 1-based sequence invariant for the AnonymousUser ID=1
- [x] Prints summary output at completion and instructions on how to start the server

## Done

- Implemented `fakedb` subcommand in `server/cmd/fakedb.go`
- Added chronological realistic logging algorithm spanning 10 years and 2000 jumps
- Grouped jumps by active days with consistent location and aircraft
- Distributed Dropzones accurately (80% local, 20% event)
- Populated disciplines with FF, WS, HOP, and FS and paired accurate altitudes and freefall times
- Scripted equipment progression downsampling canopy sizes over a skydiver's lifetime
- Successfully executes without error and correctly initializes the GORM DB with the AnonymousUser 1
