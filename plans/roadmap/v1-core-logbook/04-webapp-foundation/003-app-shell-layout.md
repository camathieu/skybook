---
ticket: "003"
epic: webapp-foundation
milestone: v1
title: App Shell Layout
status: done
priority: high
estimate: M
---

# App Shell Layout

Create the main application shell: header, navigation, content area, and responsive sidebar.

## Acceptance Criteria

- [x] `App.vue` — root layout with header bar and `<router-view>`
- [x] Header: SkyBook logo/title, jump count badge, navigation (future: tabs for BASE/Tunnel)
- [x] `+ New Jump` CTA button in the header
- [x] Responsive: hamburger menu on mobile, full nav on desktop
- [x] Keyboard shortcut hints (N for new jump, / for search)
- [x] Smooth page transitions

## Technical Notes

- Header uses gradient accent (sunset orange → teal) matching design tokens
- Jump count badge uses monospace font with counter animation (counter to be implemented in fully-functional UI tickets)

## Done

- Implemented `App.vue` app shell layout and `AppHeader.vue`.
- Added nav tab links for Jumps, and disabled links for BASE and Tunnel features.
- Implemented full responsiveness with mobile hamburger menu drawer showing sliding animation.
- Includes keyboard shortcut styling (`<kbd>`) on desktop views that is hidden on mobile views.
