---
ticket: "005"
epic: jump-api
milestone: v1
title: Delete Jump (Renumber)
status: planned
priority: high
estimate: S
---

# Delete Jump (Renumber)

`DELETE /api/v1/jumps/:id` — Delete a jump and renumber subsequent jumps downward.

## Acceptance Criteria

- [ ] Deletes the jump by ID
- [ ] All jumps with `Number > deleted.Number` shift down by 1
- [ ] Wrapped in a transaction
- [ ] Returns `204 No Content` on success
- [ ] Returns `404 Not Found` if jump doesn't exist
