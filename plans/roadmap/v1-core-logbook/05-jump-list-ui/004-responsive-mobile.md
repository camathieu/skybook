---
ticket: "004"
epic: jump-list-ui
milestone: v1
title: Responsive & Mobile
status: done
priority: medium
estimate: M
---

# Responsive & Mobile

Make the jump list work well on mobile devices.

## Acceptance Criteria

- [x] On screens < 768px: switch to card view (one card per jump)
- [x] Card shows: #, Date, Dropzone, Type, key flags
- [x] Swipe-to-delete gesture (with confirmation)
- [x] Tap card to expand details / open edit
- [x] Filter bar collapses to a "Filters" button on mobile
- [x] Touch-friendly hit targets (min 44px)

## Done

- Created `webapp/src/components/JumpCard.vue` — mobile card with 2-column grid (#, date, dropzone, type, altitude, aircraft), flag badges (Night/O₂/Cutaway), scale-on-tap animation
- `JumpList.vue` switches table → cards at `<768px` via `.desktop-only`/`.mobile-only` CSS classes
- `FilterBar.vue` collapses to "Filters" toggle button on mobile with badge count; expands to full-width stacked controls
- All interactive elements ≥ 44px touch targets; verified at 375px viewport
