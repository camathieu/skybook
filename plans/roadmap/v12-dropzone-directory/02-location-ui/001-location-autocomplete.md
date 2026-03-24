---
ticket: "001"
epic: location-ui
milestone: v12
title: Location Autocomplete
status: planned
priority: high
estimate: M
---

# Location Autocomplete

## Acceptance Criteria

- [ ] Jump form: Dropzone field queries `GET /api/v1/dropzones/autocomplete` for the shared directory
- [ ] BASE form: Location field queries `GET /api/v1/exit-points/autocomplete`
- [ ] Tunnel form: Tunnel field queries `GET /api/v1/wind-tunnels/autocomplete`
- [ ] Autocomplete results show name + country/city context
- [ ] Selecting a result sets the FK; typing a new value keeps legacy freeform string
- [ ] Fallback: legacy `GET /api/v1/jumps/autocomplete/dropzone` still works for user's own history
- [ ] API client methods in `webapp/src/api.js`
