---
milestone: v13
title: Wingloading Calculator
status: planned
---

# v13 — Wingloading Calculator

A dedicated utility in the webapp to calculate wingsuit and canopy wingloading.

## Epics

- [01 — Wingloading UI](v13-wingloading-calculator/01-wingloading-ui.md)

## Overview

### Wingloading Logic

The formula for canopy wingloading is:
`Wingloading = (Body Weight + Gear Weight) / Canopy Size`

- **Body Weight**: User's body weight (in lbs or kg).
- **Gear Weight**: Estimated or measured weight of the rig, jumpsuit, helmet, etc. (Default: 20-25 lbs).
- **Canopy Size**: The size of the canopy in sq ft (drawn from `Gear` items if v11 is completed, or manually entered).

The UI should allow swapping between Imperial (lbs) and Metric (kg), but wingloading is universally expressed in `lbs / sq ft`.

- 1 kg = 2.20462 lbs.

### Calculator UI

- A simple, standalone page or modal accessible from the Equipment/Gear section or the main menu.
- Sliders or number inputs for Body Weight and Gear Weight.
- Select from user's canopies (if v11 exists) or manual input for Canopy Size.
- Interactive output: show the current wingloading.
- Add an informational table/chart indicating experience level recommendations for specific wingloading brackets (e.g. USPA canopy sizing chart).

## Links

- PRD §10 (Future Considerations)
