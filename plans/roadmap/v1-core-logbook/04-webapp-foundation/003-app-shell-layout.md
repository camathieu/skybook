---
ticket: "003"
epic: webapp-foundation
milestone: v1
title: App Shell Layout
status: planned
priority: high
estimate: M
---

# App Shell Layout

Create the main application shell: header, navigation, content area, and responsive sidebar.

## Acceptance Criteria

- [ ] `App.vue` — root layout with header bar and `<router-view>`
- [ ] Header: SkyBook logo/title, jump count badge, navigation (future: tabs for BASE/Tunnel)
- [ ] `+ New Jump` CTA button in the header
- [ ] Responsive: hamburger menu on mobile, full nav on desktop
- [ ] Keyboard shortcut hints (N for new jump, / for search)
- [ ] Smooth page transitions

## Technical Notes

- Header uses gradient accent (sunset orange → teal) matching design tokens
- Jump count badge uses monospace font with counter animation
