---
ticket: "011"
epic: code-quality
milestone: v1.2
title: Numeric Bounds Validation
status: planned
priority: low
estimate: S
---

# Numeric Bounds Validation

Currently, numeric fields like altitude, freefall time, and canopy size have no upper bounds enforced at the API level (they just can't be negative). It is possible to submit completely unrealistic values.

## Acceptance Criteria

- [ ] Add upper bound validation to the `Jump` model (e.g., inside `Jump.Validate()`).
- [ ] Enforce `altitude` ≤ 50,000 (standard max is ~18k-30k HALO, absolute record is higher but 50k is a safe ceiling for practical data integrity).
- [ ] Enforce `freefall_time` ≤ 600 seconds (10 mins).
- [ ] Enforce `canopy_size` ≤ 500 sq ft.
- [ ] Ensure validation errors return 400 Bad Request with a clear message.
- [ ] Add backend unit tests to verify the bounds are enforced (both valid and invalid cases).
