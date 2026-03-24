---
ticket: "001"
epic: auth
milestone: v7
title: Google OAuth Setup
status: planned
priority: high
estimate: L
---

# Google OAuth Setup

## Acceptance Criteria

- [ ] Google OAuth 2.0 consent flow (`/auth/google/login`, `/auth/google/callback`)
- [ ] Config: `[auth].GoogleClientID`, `[auth].GoogleClientSecret`
- [ ] Creates user on first login (Provider="google", ProviderID=email)
- [ ] Falls back to anonymous mode when auth is not configured
