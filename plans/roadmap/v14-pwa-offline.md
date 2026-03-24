---
milestone: v14
title: PWA & Offline Access
status: planned
---

# v14 — PWA & Offline Access

Make the SkyBook webapp installable as a Progressive Web App (PWA) with offline capabilities.

## Epics

- `01-pwa-foundation` — Web manifest, service worker setup, and Vite PWA plugin integration.
- `02-offline-storage` — LocalStorage/IndexedDB caching for jumps and offline queuing for edits/creates.
- `03-sync-engine` — Background sync to push queued offline changes when the network connection is restored.
