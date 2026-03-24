---
milestone: v1
title: Core Logbook
status: done
---

# v1 — Core Logbook

> Single-user, anonymous mode. No login required.

The foundation of SkyBook: a fully functional skydive logbook with CRUD operations, auto-numbering, search/filter, and a polished dark-themed UI.

## Epics

- [Backend Foundation](v1-core-logbook/01-backend-foundation.md) — Project scaffolding, config, database, server skeleton
- [Jump Data Model](v1-core-logbook/02-jump-data-model.md) — Jump and User models, migrations, numbering invariant
- [Jump API](v1-core-logbook/03-jump-api.md) — REST endpoints for jump CRUD, insert/delete renumbering, autocomplete
- [Webapp Foundation](v1-core-logbook/04-webapp-foundation.md) — Vite + Vue + Tailwind setup, design system, SPA embedding
- [Jump List UI](v1-core-logbook/05-jump-list-ui.md) — Table component, search/filters, pagination, responsive layout
- [Jump Form UI](v1-core-logbook/06-jump-form-ui.md) — Create/edit modals, autocomplete fields, delete confirmation
- [Build & Deploy](v1-core-logbook/07-build-and-deploy.md) — Makefile, Dockerfile, AGENTS.md, ARCHITECTURE.md
- [Testing](v1-core-logbook/08-testing.md) — Backend unit tests, frontend unit tests, E2E tests
- [README](v1-core-logbook/09-readme.md) — Project README with quick start, features, and build instructions

## Goals

- A user can log, view, edit, and delete skydive jumps
- Jump numbers are always contiguous (insert/delete triggers renumbering)
- Search and filter jumps by any field
- Works on desktop and mobile with a premium dark UI
- Ships as a single Go binary with embedded SPA
