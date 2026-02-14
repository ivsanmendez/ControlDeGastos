# Memory Bank

Read files in numbered order. Lower numbers = higher priority for context loading.

## Core Documents (always read)
| File | Purpose |
|------|---------|
| `01_projectbrief.md` | Project scope, goals, and tech stack overview |
| `02_techstack.md` | Tool versions, runtime details, port mappings |
| `03_architecture.md` | Repo layout, service diagram, data flow |
| `04_activecontext.md` | Current phase, recent decisions, next steps |
| `05_progress.md` | What works, what's left, known issues |

## Subdirectories
| Directory | Purpose |
|-----------|---------|
| [`decisions/`](decisions/) | Architecture Decision Records (ADRs) |
| [`features/`](features/) | Feature-specific design docs and requirements |
| [`api/`](api/) | API endpoint docs, schemas, contracts |
| [`frontend/`](frontend/) | Component docs, state management, routing |

## Naming Convention
- Files use `NN_name.md` format (e.g., `01_projectbrief.md`)
- Lower numbers = higher priority / read first
- Leave gaps in numbering (01, 02, 05, 10...) to allow inserting later