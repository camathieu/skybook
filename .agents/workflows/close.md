---
description: Close a completed ticket (verify, review, commit, update roadmap)
---

# Close a Ticket

Finalize work on a ticket — verify completeness, run the review workflow, commit, and update all roadmap files.

## When to Use

- When you've finished implementing a ticket and want to close it out
- Invoked via `/close`

## Input

The user provides a ticket reference — either a file path, ticket number, or the currently active ticket from context.

// turbo-all

## Steps

### 1. Identify the ticket

Read the ticket file to get the full acceptance criteria:

```bash
cat <ticket-file>
```

### 2. Verify all acceptance criteria

Go through each `- [ ]` item in the ticket and verify it's actually done:

- **Code exists**: The files and functions mentioned are implemented
- **Tests pass**: Related tests exist and pass
- **Behavior works**: The feature works as described (check via tests or manual steps)

If any criterion is NOT met, stop and tell the user what's missing. Do not proceed to close an incomplete ticket.

### 3. Run `/review-changes`

Execute the full `/review-changes` workflow to catch any issues:

- Lint, build, test
- Code review checklist
- Documentation check

Fix any critical or suggested issues before proceeding.

### 4. Update the ticket file

Check off all acceptance criteria and add a `## Done` section summarizing what was delivered:

```markdown
- [x] First acceptance criterion
- [x] Second acceptance criterion

## Done

- Created `server/common/jump.go` with all v1 fields
- Added migration in `server/metadata/migrations.go`
- Unit tests in `server/common/jump_test.go` (12 tests, all passing)
```

Set the ticket status to `done`:

```yaml
status: done
```

### 5. Update parent roadmap files

Check if parent epic and milestone statuses need updating:

```bash
# Read epic to see if all sibling tickets are done
cat <epic-file>
ls <epic-dir>/
```

- If **all tickets** in the epic are now `done` → set the epic `status: done`
- If **all epics** in the milestone are now `done` → set the milestone `status: done`

### 6. Update project docs

Check if the work changed anything that affects the foundational docs:

- **`ARCHITECTURE.md`** — new packages, models, API endpoints, patterns?
- **`AGENTS.md`** — new build commands, key files, conventions?
- **`PLANS.md`** — new roadmap conventions or lifecycle changes?
- **`docs/`** — user-facing documentation for new features?

Update any affected files.

### 7. Run `/commit`

Execute the `/commit` workflow to stage, commit, and push the changes. The commit message should reference the ticket:

```
feat(<scope>): <summary>

Closes ticket <milestone>/<epic>/<ticket>
```

## Important Notes

- **Never close a ticket with unchecked acceptance criteria** — if something wasn't done, either do it now or explicitly discuss with the user why it should be dropped
- **The `## Done` section is mandatory** — it serves as the historical record of what was actually delivered
- **Roadmap cascading is important** — always check parent epic and milestone statuses
- **Docs are not optional** — per project rules, ARCHITECTURE.md and AGENTS.md must stay in sync
