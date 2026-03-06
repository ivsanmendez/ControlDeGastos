## Summary
<!-- Brief description of what this PR does and why -->

## Type of Change
<!-- Check the one that applies. Aligns with conventional commit types. -->
- [ ] `feat` — New feature
- [ ] `fix` — Bug fix
- [ ] `refactor` — Code restructuring (no behavior change)
- [ ] `docs` — Documentation only
- [ ] `test` — Adding or updating tests
- [ ] `chore` — Maintenance (deps, config, tooling)
- [ ] `ci` — CI/CD pipeline changes
- [ ] `style` — Code style (formatting, no logic change)
- [ ] `build` — Build system or external dependencies

## Scope
<!-- Check all that apply. Matches project labels. -->
- [ ] `backend` — Go API
- [ ] `frontend` — React SPA
- [ ] `infrastructure` — Docker, CI/CD, tooling
- [ ] `domain` — Business logic and domain core
- [ ] `agentic` — AI agent and event system

## Changes
<!-- List the key changes made in this PR -->
-

## Architecture Checklist
<!-- For backend changes. Skip if not applicable. -->
- [ ] Domain core has zero external dependencies
- [ ] Adapters depend on ports, not the other way around
- [ ] New interfaces defined by consumer (Go idiom)
- [ ] Only `cmd/api/main.go` imports concrete adapter types

## Testing
- [ ] Tests pass locally (`mise run test`)
- [ ] Docker build succeeds (`docker compose build`)
- [ ] Linting passes (`mise run lint`)
- [ ] Frontend tests pass (`mise run test:web`) _(if frontend changes)_
- [ ] Frontend linting passes (`mise run lint:web`) _(if frontend changes)_

## Related Issues
<!-- Link related issues. Use closing keywords to auto-close on merge. -->
<!-- Closes #issue_number -->