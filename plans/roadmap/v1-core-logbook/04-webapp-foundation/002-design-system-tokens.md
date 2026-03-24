---
ticket: "002"
epic: webapp-foundation
milestone: v1
title: Design System & Tokens
status: done
priority: high
estimate: M
---

# Design System & Tokens

Create the SkyBook design system with Tailwind theme tokens, typography, and color palette.

## Acceptance Criteria

- [x] `webapp/src/style.css` with Tailwind `@theme` customization
- [x] Color palette: deep navy/charcoal base, sunset orange → teal gradient accents
- [x] Typography: Inter (body), JetBrains Mono (numbers/monospace)
- [x] CSS custom properties for semantic colors (surface, text, accent, danger, warning)
- [x] Dark mode as default (no light mode toggle in v1)
- [x] Reusable utility classes for cards, buttons, inputs, badges

## Done

- Configured global variables and `@theme` block in `style.css` with the specified color palette and fonts.
- Extracted reusable components like `.btn-primary`, `.btn-secondary`, `.card`, and `.input` via `@layer components`.
- Styled scrollbars and defined basic vue router transition components.
