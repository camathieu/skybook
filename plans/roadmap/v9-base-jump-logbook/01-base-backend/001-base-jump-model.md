---
ticket: "001"
epic: base-backend
milestone: v9
title: BASE Jump Model
status: planned
priority: high
estimate: M
---

# BASE Jump Model

## Acceptance Criteria

- [ ] BaseJump GORM model per PRD §3.3
- [ ] Fields: Number, Date, Object, Location, Altitude, Delay, Equipment, PilotChute, Slider, WingsuitFlown, Tracking, Description, Links, Landing, CutAway, Buddies
- [ ] Independent numbering sequence (same invariant as skydive jumps)
- [ ] Migration creates `base_jumps` table with UserID FK
