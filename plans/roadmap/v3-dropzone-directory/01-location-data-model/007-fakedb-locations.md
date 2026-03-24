---
ticket: "007"
epic: location-data-model
milestone: v3
title: FakeDB integration for Locations
status: planned
priority: medium
estimate: S
---

# FakeDB integration for Locations

Extend the `skybook fakedb` algorithm to link generated jumps to the Location Database instead of raw strings.

## Acceptance Criteria

- [ ] Resolves the hardcoded string dropzones to randomly seeded `Dropzone` models in the location tables.
- [ ] Generates `WindTunnel` sessions and `ExitPoint` BASE jumps if those modules have FakeDB coverage.
- [ ] Ensures foreign keys properly link the `Jump` records to the canonical models.
