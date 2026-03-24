---
ticket: "001"
epic: tunnel-backend
milestone: v12
title: Tunnel Session Model
status: planned
priority: high
estimate: M
---

# Tunnel Session Model

## Acceptance Criteria

- [ ] TunnelSession GORM model per PRD §3.4
- [ ] Fields: Number, TotalTime (computed), Date, Tunnel, Duration, Discipline, Coach, Speed, Description, Links, Buddies
- [ ] Sequential numbering with same insert/delete renumbering as jumps
- [ ] Migration creates `tunnel_sessions` table with UserID FK
