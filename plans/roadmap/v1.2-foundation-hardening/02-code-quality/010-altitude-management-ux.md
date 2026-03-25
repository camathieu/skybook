---
ticket: "010"
epic: code-quality
milestone: v1.2
title: Altitude Management UX
status: planned
priority: medium
estimate: S
---

# Altitude Management UX

Currently, the Altitude field in `JumpModal.vue` is a standard `<input type="number">`. This requires users to manually increment/decrement or type out full numbers, which is tedious on mobile. Since most jumps occur at standard drop altitudes (e.g., 5k, 10k, 13k, 15k, 20k ft), we want to provide a faster UX.

## Acceptance Criteria

- [ ] Modify `AutocompleteInput.vue` to accept a new optional `options` prop (an array of strings).
- [ ] If `options` is provided, `AutocompleteInput.vue` should use it to populated suggestions client-side without making an API call, using standard array filtering vs the user input. => same as every other drop down in the modal sorted by last usage + possibility to add whaterver you want within the bounds (see next ticket).
- [ ] Change the Altitude field in `JumpModal.vue` to use `<AutocompleteInput id="f-altitude" :options="['5000', '10000', '12000', '13000', '15000', '20000']" />`.
- [ ] The user must still be able to input custom values (e.g., `12500`) that aren't in the static dropdown list.
- [ ] Before sending the payload to the API (`buildPayload` in `JumpModal.vue`), the string value should be parsed into a Number.
