---
ticket: "003"
epic: ux-refinements
milestone: v1.1
title: Fix Per Page
status: done
priority: high
estimate: XS
---

# Fix Per Page

Fix the bug where the "per page" pagination setting does not apply correctly.

## Acceptance Criteria

- [ ] Investigate `stores/jumps.js` and the backend `ListJumps` handler.
- [ ] Ensure modifying the `perPage` state correctly triggers an API re-fetch.
- [ ] Ensure the backend correctly applies the `limit` parameter.
- [ ] The table should accurately reflect the new page size constraint.
