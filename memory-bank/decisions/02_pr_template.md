# ADR-02: Pull Request Template

## Context
The project needed a structured PR template that enforces consistency and aligns with existing conventions: conventional commits, hexagonal architecture rules, Kanban workflow, and project labels.

## Options Considered
1. **Minimal template** — Summary + changes + testing checklist (original template)
2. **Convention-aligned template** — Sections that mirror commit types, scopes/labels, architecture rules, and issue linking

## Decision
Adopted option 2: a structured PR template at `.github/pull_request_template.md` with the following sections:

| Section | Purpose | Aligns with |
|---------|---------|-------------|
| Summary | What and why | PR description best practices |
| Type of Change | Checkbox for commit type | Conventional Commits (`feat`, `fix`, `refactor`, etc.) |
| Scope | Checkbox for affected area | Project labels (`backend`, `frontend`, `infrastructure`, `domain`, `agentic`) |
| Changes | Bullet list of key changes | Code review clarity |
| Architecture Checklist | Hexagonal rules verification | Dependency rule (domain has zero external deps, adapters depend on ports) |
| Testing | Verification steps | CI pipeline jobs (`mise run test`, `docker compose build`, `mise run lint`) |
| Related Issues | Closing keywords for auto-linking | Kanban workflow (issues move to Done on merge) |

## Consequences
- PRs are self-documenting and consistent
- Reviewers can quickly verify architectural compliance
- GitHub auto-closes linked issues on merge, keeping the Kanban board up to date
- Slight overhead for contributors filling out the template (mitigated by checkboxes and optional sections)