---
ticket: "003"
epic: tunnel-backend
milestone: v12
title: Cumulative Time Calculation
status: planned
priority: high
estimate: S
---

# Cumulative Time Calculation

## Acceptance Criteria

- [ ] `TotalTime` computed as `SUM(Duration)` of all sessions up to and including current
- [ ] Recomputed on insert/delete/update (alongside renumbering)
- [ ] Returned in API responses and displayed prominently in the UI
- [ ] Displayed in hours:minutes format (e.g. "2h 45m")
