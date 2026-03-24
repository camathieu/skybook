---
ticket: "003"
epic: jump-data-model
milestone: v1
title: Jump Number Invariant
status: planned
priority: high
estimate: M
---

# Jump Number Invariant

Implement the logic that keeps jump numbers as a contiguous 1-based sequence.

## Acceptance Criteria

- [ ] **Append**: new jump gets `Number = MAX(Number) + 1` for the user
- [ ] **Insert** at position N: all jumps with `Number >= N` shift up by 1, new jump gets N
- [ ] **Delete** jump N: remove jump, all jumps with `Number > N` shift down by 1
- [ ] All operations wrapped in a database transaction
- [ ] Unit tests for append, insert-at-start, insert-in-middle, delete-first, delete-last, delete-middle
- [ ] Edge cases: insert/delete on empty logbook, single-jump logbook

## Technical Notes

- Use `UPDATE jumps SET number = number + 1 WHERE number >= ? AND user_id = ?` in a transaction
- Consider `ORDER BY number DESC` for the shift-up to avoid unique constraint violations mid-transaction
- This logic lives in `server/metadata/` as reusable functions
