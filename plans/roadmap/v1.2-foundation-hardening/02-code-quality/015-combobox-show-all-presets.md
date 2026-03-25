---
ticket: "015"
epic: code-quality
milestone: v1.2
title: Combobox Dropdown Shows All Presets When Value Is Set
status: planned
priority: medium
estimate: S
---

# Combobox Dropdown Shows All Presets When Value Is Set

## Context

When an `AutocompleteInput` combobox already has a value (e.g. "PTU"), opening the dropdown only shows that single matching value. This is not useful — the user already knows what's filled in. The dropdown should show the full preset list (as if the input were empty) so the user can quickly switch to a different value.

## Acceptance Criteria

- [ ] When the combobox input has a value and the dropdown opens, show all preset options (not just the matching one)
- [ ] Typing should still filter the list as before
- [ ] Clearing the input should also show all presets
- [ ] Works for all combobox fields: JumpType, Altitude, Landing, Pattern, Dropzone, Aircraft
