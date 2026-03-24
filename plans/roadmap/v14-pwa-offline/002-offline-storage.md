---
ticket: "002"
epic: pwa-offline
milestone: v14
title: Offline Storage
status: planned
priority: high
estimate: L
---

# Offline Storage

Implement local caching of jump data to allow viewing and editing while offline.

## Acceptance Criteria

- [ ] Integrate a client-side database (e.g., IndexedDB via idb or localforage)
- [ ] Cache API responses (`/api/v1/jumps`, `/api/v1/config`) for offline read access
- [ ] Update frontend stores to read from the local cache when offline or as a fast initial load
- [ ] Queue local mutations (create, update, delete) in an offline outbox when the network is unavailable
- [ ] Provide a UI indicator showing when the app is "Offline" and when there are unsynced changes
