---
ticket: "005"
epic: gear-data-model
milestone: v4
title: FakeDB integration for Gear
status: planned
priority: medium
estimate: S
---

# FakeDB integration for Gear

Extend the `skybook fakedb` algorithm to populate the `Gear` and `Kit` tables and map them to jumped records, mimicking gear progression over time.

## Acceptance Criteria

- [ ] Modifies the FakeDB algorithm to create `Gear` records representing chronologically realistic canopies and containers (e.g. Nav 210 -> Sabre2 170 -> Valkyrie 103).
- [ ] Correctly links the generated jumps to `jump_gear` associations based on the active timeframe.
