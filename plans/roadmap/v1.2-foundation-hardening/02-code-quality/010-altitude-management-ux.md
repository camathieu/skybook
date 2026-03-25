---
ticket: "010"
epic: code-quality
milestone: v1.2
title: Altitude Management UX
status: done
priority: medium
estimate: S
---

# Altitude Management UX

Currently, the Altitude field in `JumpModal.vue` is a standard `<input type="number">`. This requires users to manually increment/decrement or type out full numbers, which is tedious on mobile. Since most jumps occur at standard drop altitudes (e.g., 5k, 10k, 13k, 15k, 20k ft), we want to provide a faster UX.

## Acceptance Criteria

- [x] Modify `AutocompleteInput.vue` to accept a new optional `options` prop (an array of strings).
- [x] If `options` is provided, `AutocompleteInput.vue` uses client-side filtering (no API call). On focus shows all options; on type filters by prefix.
- [x] Change the Altitude field in `JumpModal.vue` to use `AutocompleteInput` with static presets.
- [x] The user can still input custom values (e.g., `12500`) not in the preset list.
- [x] Before sending the payload, the string value is parsed into a Number (`buildPayload` already handled this).

## Done

- `webapp/src/components/AutocompleteInput.vue` — Added `options` prop (static client-side suggestions, no API call when provided) and `inputmode` prop (pass-through to `<input>` for mobile numeric keyboard). `field` stays required.
- `webapp/src/components/JumpModal.vue` — Replaced `<input type="number">` for altitude with `<AutocompleteInput field="altitude" :options="ALTITUDE_PRESETS" inputmode="numeric" />`. `ALTITUDE_PRESETS = ['5000', '10000', '12000', '13000', '15000', '20000']`.
- `webapp/src/components/AutocompleteInput.spec.js` — Added 3 tests: shows all options on focus (no API call), filters by prefix, allows custom value not in options.
- 58/58 tests pass, frontend builds cleanly.
