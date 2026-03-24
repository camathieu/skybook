---
ticket: "003"
epic: pwa-offline
milestone: v14
title: Sync Engine
status: planned
priority: high
estimate: L
---

# Sync Engine

Implement background synchronization to push queued local changes to the server when the network connection is restored.

## Acceptance Criteria

- [ ] Detect network restoration (online events)
- [ ] Process the offline outbox queue sequentially
- [ ] Handle sync conflicts or errors gracefully (e.g., keeping failed items in the queue with an error state)
- [ ] Refresh local caches with server updates after sync
- [ ] Provide UI feedback during sync ("Syncing...") and upon completion ("Up to date")
