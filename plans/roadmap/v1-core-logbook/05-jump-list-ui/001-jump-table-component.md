---
ticket: "001"
epic: jump-list-ui
milestone: v1
title: Jump Table Component
status: done
priority: high
estimate: L
---

# Jump Table Component

Core table component displaying jumps with sortable columns.

## Acceptance Criteria

- [x] Pinia store (`jumpStore`) managing jump data, loading state, and pagination
- [x] Table columns: #, Date, Dropzone, Aircraft, Type, Altitude, Freefall, Landing, Flags
- [x] Column sorting by clicking headers (toggles asc/desc)
- [x] Jump number in monospace font with accent color
- [x] Boolean flags (Night, O₂, Cutaway) shown as badges/icons
- [x] Click row to open edit modal
- [x] Smooth row entrance animations
- [x] Loading skeleton while fetching

## Done

- Created `webapp/src/stores/jumps.js` — Pinia store managing items, pagination, sort, filters, loading/error state, and bidirectional URL query sync
- Created `webapp/src/components/JumpTable.vue` — desktop table with 9 sortable columns, teal accent jump numbers (monospace), flag badges (🌙 Night, O₂, ✂ Cutaway), row entrance animations
- Created `webapp/src/components/JumpSkeleton.vue` — shimmer loading state matching table column widths
- Updated `webapp/src/views/JumpList.vue` — wires store, components, and reactive URL ↔ store sync
