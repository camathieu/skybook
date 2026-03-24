---
ticket: "007"
epic: jump-api
milestone: v1
title: Autocomplete Endpoint
status: planned
priority: medium
estimate: S
---

# Autocomplete Endpoint

`GET /api/v1/jumps/autocomplete/:field` — Return distinct values for a field, ranked by frequency.

## Acceptance Criteria

- [ ] Supported fields: `dropzone`, `aircraft`, `equipment`, `coach`, `event`
- [ ] Returns top 20 values sorted by usage count descending
- [ ] Accepts `q` query param for prefix filtering
- [ ] Returns `400 Bad Request` for unsupported field names
- [ ] Response: `[{"value": "Skydive DeLand", "count": 42}, ...]`
