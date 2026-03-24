# AGENTS.md — SkyBook

> Entry point for AI agents working on this codebase. See [ARCHITECTURE.md](ARCHITECTURE.md) for deeper technical context.

> [!CAUTION]
> ## Mandatory Review Gate — No Exceptions
>
> **NEVER perform any git or GitHub write action without explicit user approval.** This includes:
> - `git commit`, `git push`, `git push --force-with-lease`
> - Creating branches, pull requests, issues, or comments on GitHub
> - Submitting PR reviews, merging PRs
>
> **Required process for EVERY commit/push (use the `/commit` workflow):**
> 1. `git add -A && git diff --cached --stat` — show the diff summary to the user
> 2. Propose a commit message — wait for approval
> 3. `git commit` — only after user approves both diff and message
> 4. Ask before pushing — the user must explicitly say "push" or "go ahead"
>
> This applies equally to trivial one-line changes and large refactors. There are zero exceptions.

## What is SkyBook?

SkyBook is a **self-hosted skydive logbook** that lets skydivers log, search, and analyze their jump history from any device. It ships as a single Go binary (server + embedded Vue SPA) with an SQLite database — zero external dependencies, instant setup. Future versions add BASE jumping and wind tunnel tracking.

## Tech Stack

| Layer | Tech |
|-------|------|
| Backend | Go, gorilla/mux, GORM |
| Database | SQLite (WAL mode) |
| Webapp | Vue 3, Vite 7, Tailwind CSS 4, Pinia |
| Config | TOML (`skybook.cfg`) + env var overrides (`SKYBOOK_` prefix) |
| Testing | Go `testing`, Vitest, Playwright |
| Build | Makefile, Docker multi-stage |

## Repo Layout

```
skybook/
├── AGENTS.md              ← you are here
├── ARCHITECTURE.md         ← system-wide architecture
├── Makefile                ← build orchestration
├── Dockerfile
├── go.mod / go.sum
├── .agents/                ← agentic workflows (/commit, /plan, /review-changes, etc.)
├── plans/
│   ├── PRD.md              ← product requirements document
│   └── roadmap/            ← versioned roadmap (v1–v10), milestone/epic/ticket markdown files
├── server/
│   ├── main.go             ← entrypoint (cobra)
│   ├── skybook.cfg         ← default config file
│   ├── common/             ← shared types (Jump, User, Document, etc.), config
│   ├── metadata/           ← GORM backend, migrations, queries
│   ├── handlers/           ← HTTP handler functions
│   ├── middleware/          ← auth, logging, recovery, pagination
│   ├── server/             ← router setup, backend init, SPA serving
│   └── cmd/                ← cobra CLI (serve, migrate, import, export)
├── webapp/
│   ├── index.html
│   ├── vite.config.js
│   ├── package.json
│   ├── src/
│   │   ├── main.js
│   │   ├── App.vue
│   │   ├── router.js
│   │   ├── api.js           ← HTTP client for backend
│   │   ├── stores/          ← Pinia stores (jumps, base, tunnel, auth, ui)
│   │   ├── components/      ← reusable UI components
│   │   ├── views/           ← page-level components
│   │   ├── locales/         ← i18n JSON files (v8)
│   │   └── style.css        ← Tailwind entry
│   └── dist/                ← build output, embedded by Go
└── docs/                    ← GitHub Pages (VitePress)
```

## Build & Run

```bash
make all                    # Build everything (frontend + server)
make server                 # Build server only → server/skybook
make frontend               # Build Vue webapp → webapp/dist
make dev                    # Dev mode: Vite on :5173 (proxies /api to :8080) + Go on :8080
make test                   # Go unit tests
make test-frontend          # Webapp unit tests (vitest)
make test-e2e               # E2E tests (playwright)
make lint                   # go fmt + go vet
make clean                  # Remove build artifacts
make docker                 # Build Docker image
```

## Key Files

| File | Purpose |
|------|---------|
| `server/skybook.cfg` | Server configuration (TOML) — all options with comments |
| `server/common/config.go` | Config struct + parsing + env var override logic |
| `server/common/jump.go` | Jump model + type enums |
| `server/common/user.go` | User model |
| `server/metadata/jump.go` | Jump queries, numbering invariant logic |
| `server/handlers/jump.go` | Jump CRUD HTTP handlers |
| `webapp/src/api.js` | Frontend API client |
| `webapp/src/stores/jumps.js` | Jump state management (Pinia) |
| `Makefile` | Build targets for server, frontend, docker, tests |

## Conventions

- **Configuration**: TOML file + env var override using `SKYBOOK_` prefix with SCREAMING_SNAKE_CASE (e.g., `SKYBOOK_DATABASE_PATH=/data/skybook.db`)
- **Error handling**: JSON responses: `{"error": "message", "code": 400}`
- **IDs**: GORM auto-increment `uint` primary keys
- **Jump numbering invariant**: Jump numbers form a contiguous 1-based sequence per user. Insert/delete operations renumber affected jumps in a database transaction. This applies to all sequential tables (jumps, BASE jumps, tunnel sessions).
- **Multi-tenant readiness**: All user-scoped tables have a `UserID` FK from v1. In single-user mode, an anonymous user (ID=1) is auto-created.
- **API prefix**: All REST endpoints under `/api/v1/`
- **Commit messages**: [Conventional Commits](https://www.conventionalcommits.org/) with a mandatory `ticket:` trailer when the work relates to a roadmap ticket. The trailer value is the ticket's relative path, e.g. `ticket: plans/roadmap/v1.1-core-ux/01-ux-refinements/003-fix-per-page.md`

## Best Practices

- **No backward compatibility (Pre-V1)**: Until the first official release, there is zero need to maintain backward compatibility for database migrations, configuration, or APIs. Break things if it improves the architecture.
- **Always update docs**: When changing code, update `ARCHITECTURE.md`, `AGENTS.md`, and VitePress docs
- **Always update roadmaps**: When completing or starting work on a ticket/epic, update the `status` field in the corresponding roadmap files (see [PLANS.md](PLANS.md))
- **Run tests before committing**: `make lint && make test`
- **Follow the layering**: `common → metadata → middleware → handlers → cmd → server` — never import in the reverse direction
- **Dark-first UI**: The webapp uses a dark theme by default with aviation-inspired aesthetics (see ARCHITECTURE.md §Webapp)
- **Mobile-ready UI**: The mobile experience must be at least as good as desktop. Every UI component must work on a 375px screen (iPhone SE). Touch targets ≥ 44px, no hover-only interactions, tables→cards on small screens. Always verify responsiveness when implementing or reviewing UI changes.

## Workflows

This is how work gets done in this repo. All workflows are in `.agents/workflows/`.

**Standard development lifecycle:**

```
/start  →  /plan  →  (implement)  →  /review-changes  →  /close  →  /commit
```

| Workflow | Description |
|----------|-------------|
| `/start` | Start a ticket — enrich acceptance criteria, research context, create implementation plan |
| `/plan` | Create an implementation plan for a feature or change (research → design → user approval) |
| `/review-changes` | Critically review local changes (lint, build, test, full code review checklist) |
| `/close` | Close a ticket — verify all criteria, add `## Done`, cascade roadmap status, update docs |
| `/commit` | Commit and push (mandatory user review gate before any git write) |
| `/prepare-pr` | Prepare a pull request (lint, test, commit, draft PR description) |

## Documentation

1. **For agents**: `AGENTS.md` (this file) and `ARCHITECTURE.md`
2. **For humans**: VitePress site in `docs/`

## Roadmap

The full feature roadmap is in `plans/roadmap/`, organized as:
- `plans/PRD.md` — Product Requirements Document
- `plans/roadmap/<milestone>.md` — milestone overview (v1–v10)
- `plans/roadmap/<milestone>/<epic>.md` — epic breakdown
- `plans/roadmap/<milestone>/<epic>/<ticket>.md` — individual tickets with acceptance criteria
