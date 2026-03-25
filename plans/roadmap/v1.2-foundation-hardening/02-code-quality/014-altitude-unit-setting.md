---
ticket: "014"
epic: code-quality
milestone: v1.2
title: Altitude Unit Setting (Feet / Meters)
status: planned
priority: low
estimate: M
---

# Altitude Unit Setting (Feet / Meters)

## Context

SkyBook currently assumes feet for altitude values. Many European/international skydivers use meters. Add a server-side configuration setting (`altitude_unit`) and expose it to the frontend so the UI can display the correct unit label and use unit-appropriate presets.

## Acceptance Criteria

### Backend
- [ ] Add `AltitudeUnit string` to config (`skybook.cfg`) — values: `feet` (default), `meters`
- [ ] Expose the setting via `GET /api/v1/settings` (or embed in a settings response)
- [ ] No conversion logic — the server stores raw numbers, the unit is display-only

### Frontend
- [ ] Fetch the altitude unit setting on app init
- [ ] Display "Exit Altitude (ft)" or "Exit Altitude (m)" based on setting
- [ ] Use unit-appropriate presets:
  - **Feet**: `['5000', '10000', '12000', '13000', '15000', '20000']`
  - **Meters**: `['1500', '3000', '3600', '4000', '4500', '6000']`
- [ ] Update `settings.js` (from ticket 013) to hold both preset sets

### FakeDB
- [ ] Respect the altitude unit when generating data (use feet by default)
