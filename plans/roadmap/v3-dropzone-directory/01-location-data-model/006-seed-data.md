---
ticket: "006"
epic: location-data-model
milestone: v3
title: Seed Data
status: planned
priority: medium
estimate: M
---

# Seed Data

## Acceptance Criteria

- [ ] Bundled JSON/CSV seed file with common dropzones (worldwide top ~200)
- [ ] Bundled seed file with well-known BASE exit points
- [ ] Bundled seed file with known wind tunnels
- [ ] Seed logic runs on migration / first startup — only inserts if tables are empty
- [ ] Seed data is idempotent (re-running does not duplicate entries)
- [ ] Seed files stored in `server/metadata/seeds/` or similar
