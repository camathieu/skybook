---
ticket: "010"
epic: ux-refinements
milestone: v1.1
title: Component DRY Refactor
status: done
priority: high
estimate: M
---

# Component DRY Refactor

Extract shared patterns from JumpModal and ConfirmModal into reusable foundation components and global CSS utilities. This ensures strong foundations before adding more forms (BASE jumps, tunnel sessions, settings).

## Acceptance Criteria

- [x] Global form utility classes (`.form-input`, `.form-select`, `.label`, `.required`) moved from JumpModal scoped CSS to `style.css`.
- [ ] `BaseModal.vue` wrapper component extracts shared overlay/dialog pattern (fixed overlay, backdrop, transitions, click-outside-to-close, Escape key, Teleport to body).
- [ ] `JumpModal.vue` refactored to use `BaseModal.vue` — no duplicate overlay CSS remains.
- [ ] `ConfirmModal.vue` refactored to use `BaseModal.vue` — no duplicate overlay CSS remains.
- [ ] Shared button utility classes (`.btn-primary`, `.btn-secondary`, `.btn-danger`) extracted to `style.css`.
- [ ] No visual regression — all existing UI looks identical after refactor.
- [ ] All existing unit tests pass (47 tests).
- [ ] Unit tests for `BaseModal.vue` (renders slot, emits close on overlay click, emits close on Escape).
