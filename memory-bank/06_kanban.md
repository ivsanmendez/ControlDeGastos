# Kanban Board

## GitHub Project
**URL**: https://github.com/users/ivsanmendez/projects/2

## Columns
| Column | Purpose |
|--------|---------|
| Backlog | Items not yet prioritized |
| Todo | Ready to work on |
| In Progress | Currently being worked on |
| Review | Awaiting review or testing |
| Done | Completed |

## Labels
| Label | Color | Scope |
|-------|-------|-------|
| `backend` | green | Go API related |
| `frontend` | blue | React SPA related |
| `infrastructure` | yellow | Docker, CI/CD, tooling |
| `domain` | red | Business logic and domain core |
| `agentic` | purple | AI agent and event system |

## Current Board State

### Backlog
- #6 Implement AI agent driving adapter [`agentic`, `backend`]
- #7 Replace in-memory event bus with persistent broker [`agentic`, `infrastructure`]

### Todo
- #4 Build React expense management UI [`frontend`]

### In Progress
_(none)_

### Review
- #9 Rename project from ControlDeGastos to ControlDeContabilidad [`backend`, `infrastructure`] → PR #8
- SAT certificate signing + print-sign dialog [`backend`, `frontend`] — not yet a GitHub issue

### Done
- #1 Implement PostgreSQL expense repository [`backend`, `domain`]
- #2 Set up database migrations [`backend`, `infrastructure`]
- #3 Add domain unit tests with fake adapters [`backend`, `domain`]
- #5 Design and implement authentication [`backend`] — AAA framework (all 8 phases)
