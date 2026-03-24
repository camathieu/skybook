---
ticket: "001"
epic: document-backend
milestone: v6
title: Document Model
status: planned
priority: high
estimate: S
---

# Document Model

GORM model for documents per PRD §3.5.

## Acceptance Criteria

- [ ] Document struct with all fields: Name, Type, FileName, MimeType, Size, ExpiryDate
- [ ] Migration creates `documents` table with `UserID` FK
- [ ] Type enum validation: LICENSE, INSURANCE, RIG_CHECK, MEDICAL, AAD, RESERVE_REPACK, OTHER
