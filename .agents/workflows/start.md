---
description: Start working on a roadmap ticket (enrich ticket, research, create implementation plan)
---

# Start a Ticket

Begin work on a roadmap ticket — ensure the ticket is comprehensive, research context, and produce an implementation plan for user approval.

## When to Use

- When starting work on a new ticket from `plans/roadmap/`
- Before writing any code for a planned feature
- Invoked via `/start`

## Input

The user provides a ticket reference — either a file path, ticket number, or description. If ambiguous, search `plans/roadmap/` to find the matching ticket file.

// turbo-all

## Steps

### 1. Read the ticket

Read the ticket file and its parent epic and milestone for full context.

```bash
# Example: read the ticket and its parents
cat plans/roadmap/v1-core-logbook/01-backend-foundation/001-project-scaffolding.md
cat plans/roadmap/v1-core-logbook/01-backend-foundation.md
cat plans/roadmap/v1-core-logbook.md
```

### 2. Read project context

Read the foundational docs to understand current state and conventions:

```bash
cat AGENTS.md
cat ARCHITECTURE.md
cat PLANS.md
```

### 3. Enrich the ticket

Before planning implementation, critically review the ticket's acceptance criteria. A good ticket should have:

- **Clear, testable acceptance criteria** — each item should be verifiable (not vague like "works well")
- **Complete scope** — no missing edge cases, error handling, or integration points
- **Correct technical references** — field names, API paths, model references match PRD and ARCHITECTURE.md
- **Mobile responsiveness** — if the ticket involves UI, ensure criteria include mobile behavior (375px min-width, touch targets, no hover-only states)
- **Realistic estimate** — does the estimate still make sense given what you now know?

If the ticket is thin or missing criteria, **propose improvements**:

```markdown
## Suggested Additions to Acceptance Criteria

- [ ] Error response returns 404 when jump not found
- [ ] Renumbering is tested with concurrent requests
- [ ] API response includes pagination metadata
```

Present the enriched ticket to the user for approval before proceeding. Update the ticket file with approved additions.

### 4. Update ticket status

Set the ticket to `in-progress` in its YAML frontmatter:

```yaml
status: in-progress
```

If this is the first ticket in the epic, also set the epic and milestone to `in-progress`.

### 5. Research the codebase

Study existing code relevant to this ticket:

- Grep for related functions, types, routes, and tests
- Read existing implementations of similar features (follow the patterns)
- Check PRD sections referenced by the ticket
- Note any dependencies on other tickets (are prerequisites done?)

### 6. Create the implementation plan

Run the `/plan` workflow to produce a full implementation plan. The plan must:

- Reference the ticket file and its acceptance criteria
- Map each acceptance criterion to specific code changes
- Include a verification plan that proves each criterion is met
- Follow the mandatory checklist from `/plan`

### 7. Present for approval

Use `notify_user` with:
- The enriched ticket file in `PathsToReview` (if modified)
- The implementation plan in `PathsToReview`
- `BlockedOnUser: true`

Do NOT write any implementation code until the user approves the plan.

## Important Notes

- **Ticket quality is paramount**: A vague ticket leads to vague implementation. Invest time in step 3.
- **Check dependencies**: If the ticket depends on another ticket that isn't done, flag it immediately
- **One ticket at a time**: Don't combine multiple tickets into one implementation plan unless the user explicitly asks
- **PRD is the source of truth**: Cross-check acceptance criteria against the PRD for completeness
- **Do NOT auto-commit after implementation**: Once the code is written and verified (build + tests pass), ask the user to run `/review-changes` before committing. Never run `/commit` or `/close` automatically at the end of a `/start`.
