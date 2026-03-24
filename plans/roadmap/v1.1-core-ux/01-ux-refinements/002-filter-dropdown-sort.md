---
ticket: "002"
epic: ux-refinements
milestone: v1.1
title: Filter Dropdown Sort
status: planned
priority: medium
estimate: S
---

# Filter Dropdown Sort

Refine the dropdown lists in the JumpList filter toolbar so they are sorted by occurrence frequency.

## Acceptance Criteria

- [ ] Modify `JumpList.vue` or the API serving the filter options.
- [ ] Instead of sorting alphabetically or chronologically, sort the distinct values for Aircraft, Dropzone, JumpType, Event, and LO by how many jumps have that value.
- [ ] Most frequently used values should appear at the top of the select `<options>` to make filtering faster.
