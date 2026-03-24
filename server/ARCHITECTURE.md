# Architecture — Server (`server/`)

> Internals of the SkyBook Go server. For system-wide overview, data model, and API reference, see the root [ARCHITECTURE.md](../ARCHITECTURE.md).

---

## Package Structure

```
server/
├── main.go         ← entry point (calls cmd.Execute())
├── skybook.cfg     ← default TOML configuration file
├── cmd/            ← CLI commands (cobra)
├── common/         ← shared types, config, validation, helpers
├── handlers/       ← HTTP handler functions
├── metadata/       ← GORM/SQLite backend, migrations, queries
├── middleware/     ← request pipeline (recovery, logging, request ID)
└── server/         ← HTTP server setup, router, SPA embedding
```

Dependency direction: `common → metadata → middleware → handlers → cmd → server`

> [!IMPORTANT]
> Never import in the reverse direction. If `handlers` needs a type, it must be defined in `common`.

---

## `cmd/` — CLI Commands (Cobra)

| File | Command | Description |
|------|---------|-------------|
| `root.go` | `skybook` | Root command — registers `--config` flag |
| `serve.go` | `skybook serve` | Start HTTP server (default command) |
| `fakedb.go` | `skybook fakedb` | Generate a fake database with realistic jump data for testing/demos |

Config loading: `--config` flag → `./skybook.cfg` default.

> [!NOTE]
> The `fakedb` command generates jumps with controlled randomness — sequential dates (chronologically ordered, same-day allowed) and realistic discipline distribution. Used for E2E tests and development seeding.

---

## `common/` — Shared Types & Config

### Key files

| File | Content |
|------|---------|
| `jump.go` | `Jump` struct — all v1 fields, discipline enum, JSON links |
| `user.go` | `User` struct — anonymous (v1), Google OAuth (v6) |
| `date.go` | `DateOnly` type — date-only JSON serialization, stored as midnight UTC |
| `config.go` | `Config` struct — TOML parsing + `SKYBOOK_` env var override |
| `helpers.go` | `WriteJSON`, `WriteError`, `ParseUint` response helpers |

### DateOnly Type — Design Rationale

`DateOnly` wraps `time.Time` but presents as `"YYYY-MM-DD"` to the API. Stored as midnight UTC in SQLite — no hidden semantics.

```
API layer:    "2024-03-15"  ←→  JSON
Code layer:   DateOnly{time.Time}
DB layer:     2024-03-15T00:00:00Z  ←  pure midnight UTC
```

Key methods:
- `TruncateToDay()` → strips time, keeps date at midnight UTC
- `SameDay(other)` → true if same calendar day
- `AddDays(n)` → returns date shifted by n calendar days
- `UnmarshalJSON` → accepts both `YYYY-MM-DD` and RFC3339 (time stripped)

---

## `metadata/` — GORM Backend

### Initialization

- **SQLite** with WAL mode, busy timeout 5s, foreign keys ON
- **Migrations** via [gormigrate](https://github.com/go-gormigrate/gormigrate) — schema changes are versioned functions in `migrations.go`
- **In-memory databases**: tests use `:memory:` DSN via `NewBackend` for full isolation and speed

### Jump Numbering Operations

> [!CAUTION]
> **SQLite unique constraint behavior**: SQLite checks unique constraints **immediately** on each UPDATE, not at transaction commit (deferred constraints are not supported). This means all shift operations must iterate row-by-row in the correct order:
> - **Shifting UP** (insert/move-up): iterate in `DESC` order (highest number first) so `N+1` is free before `N` tries to claim it
> - **Shifting DOWN** (delete/move-down): iterate in `ASC` order (lowest number first) so `N-1` is free before `N` tries to claim it
>
> A single `UPDATE ... SET number = number + 1 WHERE number >= N` would violate the unique constraint on the first row it touches.

### Date Validation — `validateDateOrder()`

Called inside `CreateJump`, `InsertJumpAt`, `UpdateJump`, and `MoveAndUpdateJump` transactions. Enforces `date(N) ≤ date(N+1)` at the day level. Uses `skipID` to exclude the jump being moved from neighbor lookups.

Returns descriptive errors like `"date 2024-03-15 is after jump #3 (2024-03-14)"`, surfaced to the API as 400.

### Atomic Move+Update — `MoveAndUpdateJump()`

> [!IMPORTANT]
> When a jump changes **both** its position and its fields (e.g. number + date), both operations must occur in a **single transaction**. `MoveAndUpdateJump` wraps `moveJumpTx` + `updateJumpTx` in one `db.Transaction()`. If date validation fails after the move, the entire operation rolls back.
>
> Never call `MoveJump()` + `UpdateJump()` separately for the same mutation — that creates two transactions and a partial-commit risk.

### Sort & Filter Safety

- `allowedSortFields` whitelist prevents SQL injection in ORDER BY clauses
- `allowedAutocompleteFields` whitelist maps API field names → column names
- Autocomplete uses `LOWER()` on both sides for case-insensitive matching

---

## `middleware/` — Request Pipeline

Each middleware wraps `http.Handler`. Applied globally via `router.Use()`:

| File | Middleware | Purpose |
|------|-----------|---------|
| `recovery.go` | `Recovery` | Catch panics → 500 JSON error (prevents goroutine crash from killing the server) |
| `request_id.go` | `RequestID` | Generate UUID, set `X-Request-ID` header, inject into request context |
| `logging.go` | `Logging` | Log request method/path/status/duration using `slog`; wraps `ResponseWriter` to capture status code |

Execution order: Recovery → RequestID → Logging → Handler

> [!NOTE]
> `Recovery` must be outermost — if `Logging` panics, Recovery still catches it.

---

## `handlers/` — HTTP Handlers

Each handler is a closure that captures `*metadata.Backend`:

| File | Handlers | Description |
|------|----------|-------------|
| `jump.go` | `CreateJump`, `ListJumps`, `GetJump`, `UpdateJump`, `DeleteJump` | Full jump CRUD with number-based insert |
| `jump.go` | `Autocomplete` | Distinct values for form fields (dropzone, aircraft, etc.) |
| `misc.go` | `HealthHandler`, `ConfigHandler` | Health check and server config endpoint |

### Create/Insert Logic

`CreateJump` supports two modes based on the request body:
- **No `number` field** → append (standard `CreateJump`)
- **`number` field present** → insert at position (triggers shift via `InsertJumpAt`)

Date validation is handled in the metadata layer — handlers are thin.

---

## `server/` — HTTP Server & SPA Embedding

`SkyBookServer` orchestrates everything:

1. Receives `Config` + `Backend` from `cmd/serve.go`
2. Creates `mux.Router` with middleware chain and API routes
3. Mounts SPA handler as catch-all (registered **after** all API routes)
4. Starts HTTP server on configured address/port

### SPA Embedding

```go
//go:embed all:dist
var webappFS embed.FS
```

The `dist/` directory is copied from `webapp/dist/` by the Makefile before `go build`. The SPA handler:

1. Tries to serve the exact file from the embedded FS
2. If found in `/assets/` → `Cache-Control: public, max-age=31536000, immutable` (Vite hashes filenames)
3. If `index.html` → `Cache-Control: no-cache` (so new deploys are picked up immediately)
4. If file not found → serve `index.html` for client-side routing

> [!IMPORTANT]
> API routes (`/api/*`, `/health`) are registered **before** the SPA catch-all on the router. Route precedence is determined by registration order in gorilla/mux — the SPA must come last.

### Route Registration Order Trap

The autocomplete route `/jumps/autocomplete/{field}` is registered **before** `/jumps/{id:[0-9]+}` because gorilla/mux matches routes in registration order. If `{id}` came first, `autocomplete` would match as a (non-numeric) ID and return 404.

---

## Testing

All tests use **in-memory SQLite** (`:memory:`) for isolation and speed.

| Package | Test file(s) | What's tested |
|---------|-------------|--------------|
| `common` | `jump_test.go`, `config_test.go`, `user_test.go`, `date_test.go` | Model validation, config loading/env overrides, DateOnly serialization |
| `metadata` | `backend_test.go` | Jump CRUD, numbering invariant (insert/delete/move), pagination, multi-user isolation |
| `handlers` | `jump_test.go`, `misc_test.go` | All REST endpoints via `httptest`, autocomplete, health, config |
| `middleware` | `logging_test.go`, `recovery_test.go`, `request_id_test.go` | Status capture, panic recovery → JSON, UUID context injection |
| `server` | `server_test.go` | Router construction, route registration, server lifecycle |
| `cmd` | `fakedb_test.go` | FakeDB generation, jump count and numbering integrity |
