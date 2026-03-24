---
ticket: "001"
epic: pwa-offline
milestone: v14
title: PWA Foundation
status: planned
priority: high
estimate: M
---

# PWA Foundation

Set up the web manifest, service worker, and Vite PWA plugin integration to make the webapp installable.

## Acceptance Criteria

- [ ] Add `vite-plugin-pwa` to the frontend build
- [ ] Configure `manifest.webmanifest` with icons, theme color, and name
- [ ] Implement a basic service worker that caches static assets (HTML, CSS, JS, images)
- [ ] Ensure the app can load the basic shell while offline
- [ ] Add an "Install App" prompt or indicator in the UI (if supported by browser)
