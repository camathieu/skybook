---
ticket: "002"
epic: document-backend
milestone: v3
title: File Storage
status: planned
priority: high
estimate: M
---

# File Storage

Disk-based file storage for uploaded documents.

## Acceptance Criteria

- [ ] Files stored on disk at configurable path (`[storage].DocumentPath`)
- [ ] Max file size from config (`[storage].MaxDocumentSize`)
- [ ] Files named by document ID to avoid collisions
- [ ] DB stores metadata only (no blob in DB)
