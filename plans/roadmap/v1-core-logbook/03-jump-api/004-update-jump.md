---
ticket: "004"
epic: jump-api
milestone: v1
title: Update Jump
status: planned
priority: high
estimate: S
---

# Update Jump

`PUT /api/v1/jumps/:id` — Update an existing jump.

## Acceptance Criteria

- [ ] Accepts JSON body with updatable fields (all except ID, Number, UserID, CreatedAt)
- [ ] Number change triggers reorder (move jump from old position to new position)
- [ ] Validates required fields and JumpType enum
- [ ] Returns `200 OK` with updated jump object
- [ ] Returns `404 Not Found` / `400 Bad Request` as appropriate
