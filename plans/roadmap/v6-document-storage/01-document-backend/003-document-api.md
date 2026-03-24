---
ticket: "003"
epic: document-backend
milestone: v6
title: Document API
status: planned
priority: high
estimate: M
---

# Document API

REST endpoints for document CRUD.

## Acceptance Criteria

- [ ] `GET /api/v1/documents` — list with type filter
- [ ] `POST /api/v1/documents` — multipart upload
- [ ] `GET /api/v1/documents/:id` — download with proper Content-Type
- [ ] `DELETE /api/v1/documents/:id` — delete file + metadata
