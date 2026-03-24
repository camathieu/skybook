---
ticket: "003"
epic: auth
milestone: v6
title: User Isolation
status: planned
priority: high
estimate: M
---

# User Isolation

## Acceptance Criteria

- [ ] All queries scoped by `UserID` (already in place from v1 anonymous user)
- [ ] API middleware injects authenticated user into context
- [ ] Users can only see/edit their own data
- [ ] Migrate existing anonymous data to logged-in user on first login (optional)
