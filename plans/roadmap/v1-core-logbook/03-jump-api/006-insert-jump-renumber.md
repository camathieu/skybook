---
ticket: "006"
epic: jump-api
milestone: v1
title: Insert Jump (Renumber)
status: planned
priority: high
estimate: S
---

# Insert Jump (Renumber)

`POST /api/v1/jumps` with `number` field — Insert a jump at a specific position.

## Acceptance Criteria

- [ ] When `number` is provided in the request body, insert at that position
- [ ] All jumps with `Number >= requested` shift up by 1
- [ ] Wrapped in a transaction
- [ ] Validates that requested number is between 1 and `MAX(Number) + 1`
- [ ] Returns `201 Created` with the jump at the requested position
