---
ticket: "005"
epic: code-quality
milestone: v1.2
title: Deduplicate Sort Whitelist
status: planned
priority: medium
estimate: S
---

# Deduplicate Sort Whitelist

The allowed sort fields are defined in two places:
- `metadata/jump.go:19` — package-level `allowedSortFields` map
- `handlers/jump.go:55` — inline `allowedSort` map inside `ListJumps`

Both must stay in sync. They will drift.

## Acceptance Criteria

- [ ] Remove the inline `allowedSort` map from `handlers/jump.go` `ListJumps`
- [ ] Either export `allowedSortFields` from metadata, or add a `metadata.IsAllowedSortField(s string) bool` function
- [ ] Handler validates sort using the metadata layer's single source of truth
- [ ] Existing sort tests pass (valid sort, invalid sort rejection)
