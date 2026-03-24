---
ticket: "001"
epic: ux-refinements
milestone: v1.1
title: Autocomplete Recent Sort
status: planned
priority: medium
estimate: S
---

# Autocomplete Recent Sort

Refine the autocomplete API and component to sort suggestions primarily by recency of use.

## Acceptance Criteria

- [ ] Modify the `/api/v1/jumps/autocomplete/:field` endpoint.
- [ ] Currently it groups by value and sorts by `COUNT(*) DESC`.
- [ ] Change the sorting to order by the most recent `date` the value was used, or a combined heuristic (recent usage heavily weighted).
- [ ] Ensure the webapp's AutocompleteInput automatically benefits from this updated backend sorting.
