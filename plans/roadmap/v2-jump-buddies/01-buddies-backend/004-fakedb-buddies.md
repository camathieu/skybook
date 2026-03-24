---
ticket: "004"
epic: buddies-backend
milestone: v2
title: FakeDB integration for Buddies
status: planned
priority: medium
estimate: S
---

# FakeDB integration for Buddies

Extend the `skybook fakedb` algorithm to generate and assign randomized buddies to jumps.

## Acceptance Criteria

- [ ] Generates a pool of ~20 to 50 frequent jump buddies.
- [ ] In the jump assignment algorithm, assigns 1-4 buddies to RW, FF, and Wingsuit jumps based on realistic clustering.
- [ ] Updates the buddy models and `jump_buddy` join table securely during the fakedb generation process.
