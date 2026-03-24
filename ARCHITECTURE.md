# Architecture — SkyBook (System-Wide)

> System-wide architecture, package layering, data model, API reference, and build pipeline.
> See [AGENTS.md](AGENTS.md) for quick-start agent context.

---

## System Overview

SkyBook is a self-hosted skydive logbook that ships as a **single Go binary** with an embedded Vue 3 SPA.

| Component | Location | Purpose |
|-----------|----------|---------|
| **Server** (`skybook`) | `server/` | Go HTTP server — REST API, middleware chain, SQLite via GORM |
| **Webapp** | `webapp/` | Vue 3 + Vite 7 + Tailwind CSS 4 SPA, embedded in server binary |
| **Config** | `server/skybook.cfg` | TOML configuration with `SKYBOOK_` env var overrides |

### Core Abstractions

| Type | Package | Description |
|------|---------|-------------|
| `Jump` | `server/common` | Skydive jump with sequential numbering, discipline, location, equipment |
| `BaseJump` | `server/common` | BASE jump with independent numbering, object type, slider config |
| `TunnelSession` | `server/common` | Wind tunnel session with duration, cumulative time tracking |
| `JumpBuddy` | `server/common` | Shared buddy pool across all activity types (many-to-many) |
| `Document` | `server/common` | Uploaded file metadata (license, insurance, rig check) |
| `User` | `server/common` | User account — anonymous (v1) or Google OAuth (v6) |

---

## High-Level Architecture

```mermaid
graph TD
    Browser["Browser / SPA"] -->|HTTP API| Server
    CLI["skybook CLI"] -->|HTTP API| Server

    subgraph Server ["skybook Server (server/)"]
        Router["Router (gorilla/mux)"]
        MW["Middleware Chain (recovery, logging, auth)"]
        Handlers["HTTP Handlers"]
        Router --> MW --> Handlers
    end

    Handlers --> GORM["GORM"]
    GORM --> SQLite["SQLite (WAL mode)"]

    subgraph Webapp ["Embedded SPA (webapp/dist/)"]
        Vue["Vue 3 + Pinia"]
        Tailwind["Tailwind CSS 4"]
    end

    Server -->|"embed.FS"| Webapp
```

---

## Package Layering

Dependency direction flows left to right. Packages only import from packages to their left.

```
common → metadata → middleware → handlers → cmd → server
```

| Package | Location | Purpose |
|---------|----------|---------|
| `common` | `server/common/` | Shared types (Jump, User, Document), config struct, validation, `WriteJSON`/`WriteError` helpers |
| `metadata` | `server/metadata/` | GORM/SQLite backend, gormigrate migrations, queries, numbering invariant |
| `middleware` | `server/middleware/` | `Recovery`, `Logging`, `RequestID`; auth (v6), pagination |
| `handlers` | `server/handlers/` | HTTP handlers — health, config, jump CRUD, documents, stats |
| `cmd` | `server/cmd/` | Cobra CLI: `root.go` (`--config` flag, default serve), `serve.go` |
| `server` | `server/server/` | `SkyBookServer`: router, middleware chain, graceful shutdown, SPA serving |

> [!IMPORTANT]
> Never import in the reverse direction. If `handlers` needs a type, it must be defined in `common`.

---

## Request Lifecycle

```mermaid
sequenceDiagram
    participant Browser
    participant Router
    participant Middleware
    participant Handler
    participant GORM
    participant SQLite

    Browser->>Router: HTTP Request
    Router->>Middleware: Route match
    Note over Middleware: 1. Recovery<br/>2. Logging<br/>3. Auth (v6)
    Middleware->>Handler: Context populated
    Handler->>GORM: Query / Mutate
    GORM->>SQLite: SQL
    SQLite-->>GORM: Result
    GORM-->>Handler: Models
    Handler-->>Browser: JSON Response
```

---

## Data Model

### Jump (v1)

The core entity. All fields are defined in [PRD §3.1](plans/PRD.md).

| Key Field | Type | Notes |
|-----------|------|-------|
| `ID` | `uint` (PK) | Auto-increment |
| `Number` | `uint` | Sequential, contiguous, per-user. **Subject to renumbering.** |
| `UserID` | `uint` (FK) | Multi-tenant readiness — defaults to anonymous user (ID=1) in v1 |
| `Date` | `datetime` | Required |
| `Dropzone` | `string` | Required, autocomplete from history |
| `JumpType` | `string` | Enum: FF, FS, CRW, HOP, CF, AFF, AFFI, CAMERA, TANDEM, DEMO, XRW, ANGLE, TRACKING, CP, WINGSUIT, OTHER |
| `Links` | `JSON text` | Array of URLs, stored as JSON |
| `Buddies` | `[]JumpBuddy` | Many-to-many via `jump_buddies` join table (v4) |

### BaseJump (v9) & TunnelSession (v10)

Separate tables with independent numbering sequences. Same renumbering invariant applies. Full schemas in [PRD §3.3–3.4](plans/PRD.md).

### JumpBuddy (v4)

Shared pool across all activity types. Linked via join tables: `jump_buddies`, `base_jump_buddies`, `tunnel_session_buddies`. Autocomplete ranked by popularity (total jump count across all types).

### Document (v3)

File metadata stored in DB, file content on disk at configurable path. Categories: LICENSE, INSURANCE, RIG_CHECK, MEDICAL, AAD, RESERVE_REPACK, OTHER.

### User (v1 anonymous → v6 multi-tenant)

In v1, an anonymous user (ID=1) is auto-created. All data is attributed to this user. When v6 adds Google OAuth, no schema changes are needed.

---

## Jump Number Invariant

> [!CAUTION]
> This is the most critical business rule in the system. All sequential tables (jumps, BASE jumps, tunnel sessions) follow this invariant.

**Rule**: Jump numbers form a **contiguous 1-based sequence** per user. No gaps, no duplicates.

| Operation | Behavior |
|-----------|----------|
| **Append** | `Number = MAX(Number) + 1` |
| **Insert at position N** | Shift `[N…MAX]` up by 1 (DESC order), then insert at N |
| **Delete jump N** | Remove jump, shift `[N+1…MAX]` down by 1 (ASC order) |
| **Move jump N → M** | Park at sentinel, shift intermediate range, place at M |
| **Bulk import** | Disable auto-numbering, assign final numbers, validate contiguity |

All operations are wrapped in a **database transaction** for atomicity.

> [!NOTE]
> SQLite checks unique constraints **immediately** (not deferred). All shift operations must iterate row-by-row in the correct order: **DESC** when shifting numbers up, **ASC** when shifting numbers down. The single-statement `UPDATE ... ORDER BY` would require deferred constraints which SQLite does not support.

```go
// Pseudocode — Insert at position (actual impl uses iterative updates)
func InsertJumpAt(tx *gorm.DB, userID, position uint, jump *Jump) error {
    // Load rows to shift in DESC order, update one-by-one
    var toShift []*Jump
    tx.Where("user_id = ? AND number >= ?", userID, position).
        Order("number DESC").Find(&toShift)
    for _, j := range toShift {
        tx.Model(j).Update("number", j.Number+1)
    }
    jump.Number = position
    return tx.Create(jump).Error
}
```


---

## API Reference (v1)

All endpoints prefixed with `/api/v1`. Responses use JSON.

### Jumps

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/jumps` | List jumps (paginated, filterable, sortable) |
| `POST` | `/api/v1/jumps` | Create jump (append, or insert at position if `number` specified) |
| `GET` | `/api/v1/jumps/:id` | Get single jump |
| `PUT` | `/api/v1/jumps/:id` | Update jump |
| `DELETE` | `/api/v1/jumps/:id` | Delete jump (triggers renumber) |
| `GET` | `/api/v1/jumps/autocomplete/:field` | Distinct values for field; `?sort=alpha` for A–Z (filter dropdowns), default is recency (modal autocomplete) |

### Other Endpoints

| Method | Path | Version | Description |
|--------|------|---------|-------------|
| `GET` | `/api/v1/stats` | v2 | Logbook statistics |
| `GET/POST/DELETE` | `/api/v1/documents[/:id]` | v3 | Document CRUD |
| `GET` | `/api/v1/buddies` | v4 | Buddy list with popularity |
| `POST` | `/api/v1/export` | v5 | Export logbook as JSON |
| `POST` | `/api/v1/import` | v5 | Import logbook from JSON |
| `GET` | `/api/v1/config` | v1 | Server config / feature flags |
| `GET` | `/health` | v1 | Health check |

### Query Parameters (`GET /api/v1/jumps`)

| Param | Type | Description |
|-------|------|-------------|
| `page` | int | Page number (1-based, default 1) |
| `per_page` | int | Results per page (default 25, max 100) |
| `sort` | string | `number`, `date`, `dropzone`, `altitude` |
| `order` | string | `asc` / `desc` (default: `desc`) |
| `q` | string | Full-text search (description, dropzone, event, lo) |
| `date_from` / `date_to` | date | Date range filter |
| `dropzone` | string | Exact dropzone filter |
| `jump_type` | string | Discipline filter |
| `altitude_min` / `altitude_max` | int | Altitude range |
| `cutaway` / `night` | bool | Boolean flag filters |
| `lo` | string | Load Organizer / Coach filter |

### Error Responses

All errors return a consistent JSON format:

```json
{
    "error": "descriptive error message",
    "code": 400
}
```

---

## Webapp Architecture

### Stack

- **Vue 3** with Composition API and `<script setup>`
- **Pinia** for state management (stores: `jumps`, `toast`, `base`, `tunnel`, `auth`, `ui`)
- **Vue Router** for client-side routing
- **Tailwind CSS 4** with `@theme` customization

### Design System

- **Dark-first**: Deep navy/charcoal base (`#0f172a` → `#1e293b`)
- **Accent gradient**: Sunset orange (`#f97316`) → Teal (`#14b8a6`)
- **Typography**: Inter (body), JetBrains Mono (numbers/monospace)
- **Micro-animations**: Row insertions, number counters, filter transitions

### SPA Routing

| Route | View | Description |
|-------|------|-------------|
| `/` | JumpList | Main logbook table (default) |
| `/stats` | Statistics | Dashboard (v2+) |
| `/documents` | Documents | Document storage (v3+) |
| `/base` | BaseJumpList | BASE logbook (v9) |
| `/tunnel` | TunnelList | Tunnel sessions (v10) |
| `/settings` | Settings | User preferences (v6+) |

### Jump List Components

The main logbook view follows a component hierarchy:

| Component | File | Description |
|-----------|------|-------------|
| `JumpList.vue` | `views/` | Page-level view — URL sync, state management, layout switching |
| `JumpTable.vue` | `components/` | Desktop table — sortable column headers, flag badges, row animations; emits `edit` on row click |
| `JumpCard.vue` | `components/` | Mobile card — 2-column grid layout, touch-friendly; emits `edit` on card click |
| `JumpSkeleton.vue` | `components/` | Loading shimmer — matches table dimensions |
| `SearchBar.vue` | `components/` | Debounced search (300ms) with `/` keyboard shortcut |
| `FilterBar.vue` | `components/` | Type, dropzone, date range, boolean toggles; collapses on mobile |
| `Pagination.vue` | `components/` | Prev/Next, page indicator, per-page selector (25/50/100) |

**State management**: `stores/jumps.js` (Pinia) manages items, pagination, sort, and filters — plus `createJump`, `updateJump`, `deleteJump` mutation actions. URL query params are synced bidirectionally — on mount, query → store; on change, store → URL via `router.replace`.

### Jump Form Components

The jump create/edit workflow uses three new components:

| Component | File | Description |
|-----------|------|-------------|
| `BaseModal.vue` | `components/` | Shared modal wrapper — Teleport, backdrop overlay, click-outside-to-close, Escape key, z-index prop. Used by JumpModal and ConfirmModal. |
| `JumpModal.vue` | `components/` | Unified create/edit form — pass `jump` prop for edit mode, null for create. Full-screen sheet on mobile `<640px`. Uses `BaseModal`. |
| `AutocompleteInput.vue` | `components/` | Debounced (200ms) autocomplete backed by `/api/v1/jumps/autocomplete/:field`; shows suggestions on focus and on input; keyboard navigable dropdown |
| `ConfirmModal.vue` | `components/` | Generic danger confirmation dialog with loading state; reusable for any destructive action. Uses `BaseModal`. |

**Modal trigger flow**: `JumpList` maintains `showModal: ref(bool)` and `editingJump: ref(jump|null)`. Clicking `+ New Jump` (or pressing `N`) opens create mode. Clicking a table row or card opens edit mode with the jump pre-populated. The `N` shortcut is registered globally on `window` in `onMounted` and cleaned up in `onUnmounted`.

**Mutation actions** in `stores/jumps.js`: `createJump`, `updateJump`, `deleteJump` — all call the API then trigger a `fetchJumps()` refresh to keep the list consistent with auto-renumbering.

### API Client

`webapp/src/api.js` — centralized HTTP client using `fetch()`:
- Base URL: `/api/v1`
- JSON request/response helpers
- Error handling with toast notifications
- In dev mode, Vite proxies `/api/*` to `localhost:8080`

---

## Configuration

TOML config file (`server/skybook.cfg`) with env var overrides using `SKYBOOK_` prefix:

```toml
[server]
ListenAddress = "0.0.0.0"
ListenPort = 8080
Debug = false

[database]
Path = "./skybook.db"

[auth]
Provider = "anonymous"              # "anonymous" (v1) or "google" (v6)
# GoogleClientID = ""
# GoogleClientSecret = ""

[storage]
MaxDocumentSize = "10MB"
DocumentPath = "./documents"

[defaults]
UnitSystem = "imperial"             # "imperial" or "metric"
DefaultJumpType = "FF"
```

**Env var override pattern**: `SKYBOOK_SECTION_KEY` (e.g., `SKYBOOK_DATABASE_PATH=/data/skybook.db`)

---

## Build Pipeline

### Makefile Targets

```
make all
  ├── make frontend      → cd webapp && npm ci && npm run build → webapp/dist/
  └── make server        → cp webapp/dist → server/server/dist, go build (CGO, -tags osusergo,netgo,sqlite_omit_load_extension, -static) → server/skybook

make dev
  ├── webapp: vite dev server on :5173 (proxies /api/* to :8080)
  └── server: go run on :8080
```

### SPA Embedding

```go
//go:embed all:dist
var webappFS embed.FS

func serveSPA(router *mux.Router) {
    distFS, _ := fs.Sub(webappFS, "dist")
    fileServer := http.FileServer(http.FS(distFS))
    router.PathPrefix("/").Handler(spaHandler(fileServer))
}
```

- API routes (`/api/*`, `/health`) take precedence over the SPA catch-all
- All non-matching routes fall through to `index.html` for client-side routing
- Hashed static assets get long-lived cache headers; `index.html` gets `no-cache`

### Docker

Multi-stage build:
1. **Node stage**: Build frontend (`npm ci && npm run build`)
2. **Go stage**: Build server with embedded frontend
3. **Runtime stage**: Alpine with the single binary

---

## Testing

### Backend (`make test`)

All backend packages have unit test coverage using Go's standard `testing` library and `net/http/httptest`.

| Package | Test file(s) | What's tested |
|---------|-------------|--------------|
| `common` | `jump_test.go`, `config_test.go`, `user_test.go` | Model validation, config loading/env overrides, anonymous user |
| `metadata` | `backend_test.go` | Jump CRUD, numbering invariant (insert/delete/move), pagination, multi-user isolation |
| `handlers` | `jump_test.go`, `misc_test.go` | All REST endpoints via httptest, autocomplete, health, config |
| `middleware` | `logging_test.go`, `recovery_test.go`, `request_id_test.go` | Status capture, panic recovery → 500 JSON, UUID context injection |
| `server` | `server_test.go` | Router construction, route registration, server start/shutdown lifecycle |
| `cmd` | `fakedb_test.go` | FakeDB generation, jump count and numbering integrity |

Tests use **in-memory SQLite** databases (`:memory:` via `metadata.NewBackend`) for isolation and speed.

### Frontend Unit Tests (`make test-frontend`)

Uses **Vitest** with `@vue/test-utils` and **jsdom** for fast, headless component testing.

| File | What's tested |
|------|--------------|
| `src/stores/jumps.spec.js` | State init, CRUD actions, sorting, filtering, pagination, URL sync |
| `src/components/AutocompleteInput.spec.js` | Rendering, v-model, debounced API calls, keyboard nav |
| `src/components/Pagination.spec.js` | Page range display, Prev/Next disabled states, per-page switching |

Vitest is scoped to `src/**/*.spec.js` (configured in `vite.config.js`) to avoid picking up Playwright E2E files.

### E2E Tests (`make test-e2e`)

Uses **Playwright** running against a live full stack (Go backend + Vite dev server).

The `playwright.config.js` `webServer` block auto-boots:
1. Go backend on `:8080` with `SKYBOOK_DATABASE_PATH=:memory:` for isolation
2. Vite dev server on `:5173` (proxies `/api` → `:8080`)

| File | What's tested |
|------|-------------|
| `e2e/jumps.spec.js` | Jump CRUD (create, edit, delete, insert-at, search/filter) |
| `e2e/mobile.spec.js` | Responsive layout, touch targets ≥44px, modal usability at 375px |

Projects: **Desktop Chrome** and **Mobile Safari (iPhone SE)**. All interactable elements have `data-testid` attributes.

---

## Multi-Tenant Readiness

Even in v1 (anonymous single-user mode):

- All user-scoped tables (`jumps`, `documents`, etc.) have a `UserID` FK
- A default anonymous user (ID=1, Provider="local") is auto-created on first startup
- All data is attributed to this anonymous user
- When v7 adds Google OAuth: **zero schema changes needed**, just add the auth middleware

---

## Versioned Feature Roadmap

| Version | Feature | Status |
|---------|---------|--------|
| **v1** | Core Logbook — CRUD, auto-numbering, search, dark UI | Planned |
| **v2** | Basic Statistics — charts, metrics | Planned |
| **v3** | Document Storage — upload, expiry tracking | Planned |
| **v4** | Jump Buddies — shared pool, popularity ranking | Planned |
| **v5** | Import/Export — JSON with schema versioning | Planned |
| **v6** | Multi-Tenant — Google OAuth, user isolation | Planned |
| **v7** | Advanced Statistics — heatmap, currency, records | Planned |
| **v8** | Internationalization — vue-i18n, unit system | Planned |
| **v9** | BASE Jump Logbook — separate tab, BASE-specific fields | Planned |
| **v10** | Tunnel Time Tracker — session tracking, cumulative time | Planned |
| **v11** | Gear & Kit Tracking — detailed equipment items, sizes, kits | Planned |
| **v12** | Location Directory — shared Dropzone, ExitPoint, WindTunnel tables | Planned |
| **v13** | Wingloading Calculator — standalone UI utility | Planned |

Full details: [plans/PRD.md](plans/PRD.md) and [plans/roadmap/](plans/roadmap/)
