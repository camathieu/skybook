---
ticket: "001"
epic: stats-api
milestone: v2
title: Statistics Endpoint
status: planned
priority: high
estimate: M
---

# Statistics Endpoint

`GET /api/v1/stats` — Compute and return logbook statistics.

## Acceptance Criteria

- [ ] Total jumps, total freefall time (seconds)
- [ ] Jumps per discipline (map: type → count)
- [ ] Jumps per dropzone (top 10 by count)
- [ ] Jumps per month over the last 2 years
- [ ] Average and max exit altitude
- [ ] Cutaway count
- [ ] Recent activity streak (consecutive weeks with at least 1 jump)
- [ ] Computed via SQL aggregation queries (not in-memory)
