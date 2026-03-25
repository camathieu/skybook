---
ticket: "005"
epic: code-quality
milestone: v1.2
title: Deduplicate Sort Whitelist
status: done
priority: medium
estimate: S
---

# Deduplicate Sort Whitelist

The allowed sort fields are defined in two places:
- `metadata/jump.go:19` — package-level `allowedSortFields` map
- `handlers/jump.go:55` — inline `allowedSort` map inside `ListJumps`

Both must stay in sync. They will drift.

## Acceptance Criteria

- [x] Remove the inline `allowedSort` map from `handlers/jump.go` `ListJumps`
- [x] Either export `allowedSortFields` from metadata, or add a `metadata.IsAllowedSortField(s string) bool` function
- [x] Handler validates sort using the metadata layer's single source of truth
- [x] Existing sort tests pass (valid sort, invalid sort rejection)

## Done

- Added `IsAllowedSortField(field string) bool` in `server/metadata/jump.go`; `allowedSortFields` stays unexported
- Removed inline `allowedSort` map from `handlers/jump.go`; handler now calls `metadata.IsAllowedSortField(sortBy)`
- All 6 test packages pass; `TestListJumps_InvalidSort` confirms rejection still works
