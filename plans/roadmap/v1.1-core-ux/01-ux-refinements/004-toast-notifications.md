---
ticket: "004"
epic: ux-refinements
milestone: v1.1
title: Toast Notifications
status: done
priority: medium
estimate: M
---

# Toast Notifications

A unified snackbar/toast system for non-blocking success and error feedback.

## Acceptance Criteria

- [x] Global toast notification manager (Pinia `useToastStore`) with `addToast(message, type, duration?)`.
- [x] Toast component (`ToastContainer.vue`) renders a stack of self-dismissing messages (default 4s).
- [x] **Success** styling (teal accent) for jump create, edit, and delete.
- [x] **Error** styling (red/error accent) for API failures in JumpModal (replaces inline `saveError`).
- [x] Stackable with max 5 visible; oldest auto-removed when limit exceeded.
- [x] Each toast has a close × button for manual dismiss.
- [x] Smooth slide-in / fade-out animation (CSS transitions).
- [x] Mobile: toasts render full-width, bottom-positioned, with ≥44px touch target on dismiss button.
- [x] JumpList.vue wires `@close` event from JumpModal to `addToast()` for create/edit/delete success.
- [x] Unit tests for the toast store (add, auto-remove, max cap).

## Done

- `webapp/src/stores/toast.js` (new): Pinia store with `addToast`, `removeToast`, auto-dismiss via setTimeout, max 5 cap
- `webapp/src/components/ToastContainer.vue` (new): fixed-position toast stack with slide-in/fade-out animation, success/error/info variants, aria-live, 44px mobile touch targets
- `webapp/src/stores/toast.spec.js` (new): 5 unit tests covering add, auto-remove, manual remove, max cap, defaults
- `webapp/src/App.vue`: Added `<ToastContainer />` at root level
- `webapp/src/views/JumpList.vue`: `onModalClose` fires "Jump created" / "Jump updated" / "Jump deleted" toasts
- `webapp/src/style.css`: Added `--color-info` design token
