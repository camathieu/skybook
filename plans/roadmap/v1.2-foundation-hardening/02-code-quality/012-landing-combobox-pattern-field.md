---
ticket: "012"
epic: code-quality
milestone: v1.2
title: Landing Combobox, Pattern Field & FakeDB Updates
status: done
priority: medium
estimate: M
---

# Landing Combobox, Pattern Field & FakeDB Updates

## Acceptance Criteria

### Backend
- [x] Add `Pattern string` field to `common.Jump` struct (`gorm:"size:32" json:"pattern,omitempty"`)
- [x] Add `"pattern"` to `updatableColumns` in `metadata/jump.go`
- [x] ~~Add migration~~ — No migration needed pre-v1; base `create_jumps` migration runs `AutoMigrate(&common.Jump{})` which creates all fields

### Frontend — Landing Combobox
- [x] Replace `<select>` for Landing with `<AutocompleteInput :options="LANDING_OPTIONS">` (same combobox UX as altitude)
- [x] Users can type custom landing values not in the preset list (e.g. "tree")
- [x] Preset options: `['Stand-up', 'Sliding', 'PLF', 'Off-DZ', 'Water']`

### Frontend — Pattern Field
- [x] Add Pattern combobox field after Landing in the Details fieldset
- [x] Preset options: `['PTU', 'Straight Final', '90°', '270°', '450°', '630°']`
- [x] Add `pattern` to form reactive object, `buildPayload`, and `resetForm`

### FakeDB
- [x] Populate Landing: 90% Stand-up, 9% Sliding, 1% Off-DZ
- [x] Populate Pattern with canopy progression: ≤200→PTU, 201–500→90°, 501+→270°
- [x] Regenerate `skybook.db` with 2000 jumps

### Documentation
- [x] Update `webapp/ARCHITECTURE.md` to reflect Landing combobox and new Pattern field

## Done

### Ticket scope
- Converted Landing from `<select>` to `AutocompleteInput` combobox (custom values supported)
- Added Pattern combobox field with canopy approach presets
- Converted JumpType from `<select>` to `AutocompleteInput` combobox
- Rearranged Details section into 3 paired two-column rows: JumpType+Altitude, Freefall+Canopy, Pattern+Landing
- FakeDB generates realistic distributions (landing 90/9/1%, pattern progression PTU→90°→270°)
- Deleted unnecessary extra migrations (`add_packjob`, `add_pattern`) — base schema handles all fields

### Bonus scope (user-requested during implementation)
- Added `Favorite bool` field (indexed) to Jump struct with full-stack support: star toggle in modal header, ★ filter button in JumpList toolbar, star icons in table/card rows, `?favorite=true` API filter, ~5% fakedb population
- Created roadmap tickets 013 (Server-Driven Frontend Settings) and 014 (Altitude Unit Setting)
