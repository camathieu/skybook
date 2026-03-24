---
ticket: "003"
epic: location-ui
milestone: v3
title: Activity Location Resolver
status: planned
priority: medium
estimate: M
---

# Activity Location Resolver

## Acceptance Criteria

- [ ] On jump/BASE/tunnel list views, show a "Link" action for rows with a freeform string but no FK
- [ ] "Link" action opens a modal with autocomplete to resolve the string against the canonical directory
- [ ] Fuzzy match suggestions based on the existing string value
- [ ] User can confirm a match or create a new directory entry from the string
- [ ] Resolving sets the FK and optionally clears the freeform string
- [ ] Bulk resolve option: resolve all unlinked rows matching a given string at once
