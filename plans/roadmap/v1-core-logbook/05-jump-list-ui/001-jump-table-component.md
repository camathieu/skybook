---
ticket: "001"
epic: jump-list-ui
milestone: v1
title: Jump Table Component
status: planned
priority: high
estimate: L
---

# Jump Table Component

Core table component displaying jumps with sortable columns.

## Acceptance Criteria

- [ ] Pinia store (`jumpStore`) managing jump data, loading state, and pagination
- [ ] Table columns: #, Date, Dropzone, Aircraft, Type, Altitude, Freefall, Landing, Flags
- [ ] Column sorting by clicking headers (toggles asc/desc)
- [ ] Jump number in monospace font with accent color
- [ ] Boolean flags (Night, O₂, Cutaway) shown as badges/icons
- [ ] Click row to open edit modal
- [ ] Smooth row entrance animations
- [ ] Loading skeleton while fetching
