---
ticket: "003"
epic: build-and-deploy
milestone: v1
title: AGENTS.md, ARCHITECTURE.md & Project Docs
status: done
priority: high
estimate: M
---

# AGENTS.md, ARCHITECTURE.md & Project Docs

Project documentation for AI coding assistants, architectural reference, and agentic workflows.

## Acceptance Criteria

- [x] `AGENTS.md` — AI assistant guidelines, project conventions, build commands
- [x] `ARCHITECTURE.md` — system overview, package layering, data model, API reference
- [x] `PLANS.md` — roadmap structure, naming conventions, ticket lifecycle, completion reporting
- [x] `.agents/workflows/` — port Plik workflows to SkyBook (adapt or remove Plik-specific ones)
- [x] Both docs and workflows kept in sync with actual codebase as per project rules

## Done

- Created `AGENTS.md` (134 lines) — tech stack, repo layout, build commands, conventions, best practices
- Created `ARCHITECTURE.md` (290 lines) — system overview with mermaid diagrams, package layering, data model, API reference, jump numbering invariant, webapp architecture, config, build pipeline
- Created `PLANS.md` (90 lines) — roadmap structure, naming conventions, YAML frontmatter schema, ticket lifecycle with completion reporting (`## Done` sections)
- Ported `.agents/workflows/` — kept `commit.md`, `plan.md` (generic), adapted `review-changes.md` and `prepare-pr.md`, removed `cut-release.md`, `add-language.md`, `review-language.md` (Plik-specific)
