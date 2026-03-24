---
ticket: "002"
epic: jump-data-model
milestone: v1
title: User Model (Anonymous)
status: done
priority: high
estimate: S
---

# User Model (Anonymous)

Create the User model and auto-provision an anonymous default user for v1 single-user mode.

## Acceptance Criteria

- [x] `server/common/user.go` — User struct with all fields, GORM tags, JSON serialization
- [x] Full field table:

| Field | Go Type | GORM / Notes |
|-------|---------|-------------|
| `ID` | `uint` | PK, auto-increment |
| `Provider` | `string` | Auth provider: `local` (v1), `google` (v6) |
| `ProviderID` | `string` | Provider-specific user ID (empty for local) |
| `Email` | `string` | User email (empty for anonymous) |
| `Name` | `string` | Display name — default "Skydiver" for anonymous |
| `Locale` | `string` | Preferred locale — default "en" (used by i18n in v8) |
| `UnitSystem` | `string` | `imperial` or `metric` — default "imperial" |
| `CreatedAt` | `time.Time` | auto |
| `UpdatedAt` | `time.Time` | auto |

- [x] Migration creates `users` table
- [x] On first startup, auto-create anonymous user (ID=1, Provider=`local`, Name=`Skydiver`)
- [x] All jumps in v1 are attributed to this anonymous user
- [x] Indexes: `(provider, provider_id)` unique
- [x] Unit tests for anonymous user creation and defaults

## Technical Notes

- This establishes **multi-tenant readiness**: when v6 adds Google OAuth, no schema changes needed — just add new users with Provider=`google`
- The anonymous user is created in the gormigrate initial migration (not server startup), so it exists before any handler runs
- `Locale` is stored but unused until v8 (i18n)
- `UnitSystem` is user-level preference that the frontend reads from `GET /api/v1/config` or a future `GET /api/v1/me` endpoint

## Done
- Created `server/common/user.go` and `AnonymousUser()` factory.
- Added `202603241300_create_users` migration seeded with anonymous user logic.
- Migration runs first so foreign key constraints on `Jump.UserID` are satisfied.
- Unit tests in `server/common/user_test.go` and integration test `TestAnonymousUserSeeded`.
