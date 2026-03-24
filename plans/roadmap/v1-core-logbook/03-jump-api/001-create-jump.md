---
ticket: "001"
epic: jump-api
milestone: v1
title: Create Jump
status: planned
priority: high
estimate: M
---

# Create Jump

`POST /api/v1/jumps` — Create a new jump, appended at the end by default.

## Acceptance Criteria

- [ ] Accepts JSON body with all Jump fields (except ID, Number, CreatedAt, UpdatedAt)
- [ ] Auto-assigns `Number = MAX(Number) + 1`
- [ ] Validates required fields: Date, Dropzone, JumpType
- [ ] Validates JumpType against allowed enum values
- [ ] Returns `201 Created` with the full jump object (including assigned Number)
- [ ] Returns `400 Bad Request` with descriptive error for invalid input

## Technical Notes

- Number assignment uses the invariant logic from ticket `jump-data-model/003`
- UserID is hardcoded to anonymous user (ID=1) in v1
