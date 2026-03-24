# Plans Directory

> Structure and conventions for the SkyBook planning documents.

## Contents

| Path | Purpose |
|------|---------|
| `PRD.md` | Product Requirements Document — full feature spec, data model, API design |
| `roadmap/` | Versioned roadmap (v1–v10) organized as milestone → epic → ticket |

## Roadmap Structure

```
roadmap/
├── <milestone>.md                     ← milestone overview (e.g. v1-core-logbook.md)
├── <milestone>/
│   ├── <NN>-<epic>.md                 ← epic summary, ordered by dependency
│   └── <NN>-<epic>/
│       └── <NNN>-<ticket>.md          ← individual ticket with acceptance criteria
```

### Naming Convention

- **Milestones**: `v<N>-<slug>.md` (e.g. `v1-core-logbook.md`)
- **Epics**: `<NN>-<slug>.md` with matching directory (e.g. `01-backend-foundation.md` + `01-backend-foundation/`)
- **Tickets**: `<NNN>-<slug>.md` (e.g. `001-project-scaffolding.md`)

### YAML Frontmatter

Every file has YAML frontmatter for machine-readable metadata:

```yaml
---
milestone: v1
epic: backend-foundation      # epics and tickets only
ticket: "001"                  # tickets only
title: Project Scaffolding
status: planned                # planned | in-progress | done
priority: high                 # tickets only: high | medium | low
estimate: S                    # tickets only: XS | S | M | L | XL
---
```

## Ticket Lifecycle

Every ticket follows this lifecycle:

```
planned → in-progress → done
```

### Starting Work

1. Set ticket `status: in-progress` in frontmatter
2. If this is the first ticket in the epic, set the epic `status: in-progress` too
3. If this is the first epic in the milestone, set the milestone `status: in-progress` too

### During Work

- Check off `- [ ]` acceptance criteria as they are met (`- [x]`)
- Keep ticket body up to date with any deviations or decisions

### Completing Work

1. Verify all acceptance criteria are checked off
2. Add a `## Done` section at the bottom of the ticket summarizing what was delivered:
   ```markdown
   ## Done
   
   - Created `server/common/jump.go` with all v1 fields
   - Migration tested with empty and seeded databases
   - Unit tests added in `server/common/jump_test.go`
   ```
3. Set ticket `status: done`
4. If all tickets in the epic are done, set the epic `status: done`
5. If all epics in the milestone are done, set the milestone `status: done`

## Critical Rule

> [!CAUTION]
> **Always update roadmap files when working on the codebase.**
>
> When starting, completing, or modifying work related to a ticket or epic:
> 1. Update the `status` field in the ticket's frontmatter (`planned` → `in-progress` → `done`)
> 2. Check off acceptance criteria in the ticket body as they are met
> 3. Add a `## Done` section when completing the ticket
> 4. Update parent epic and milestone statuses as appropriate
>
> Do not consider any task complete until the corresponding roadmap files are verified and updated.
