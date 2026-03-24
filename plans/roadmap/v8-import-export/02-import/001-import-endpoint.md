---
ticket: "001"
epic: import
milestone: v8
title: Import Endpoint
status: planned
priority: high
estimate: L
---

# Import Endpoint

## Acceptance Criteria

- [ ] `POST /api/v1/import` — accept JSON logbook file
- [ ] Schema version validation
- [ ] Conflict resolution modes: merge, overwrite, skip (via query param)
- [ ] Progress reporting for large imports
- [ ] Transactional — rollback on failure
