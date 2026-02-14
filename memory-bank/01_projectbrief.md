# Project Brief: ControlDeGastos

## Overview
ControlDeGastos is a full-stack expense tracking application.

## Tech Stack
- **Backend**: Go 1.23 (REST API)
- **Frontend**: React (Vite + TypeScript) — SPA
- **Database**: PostgreSQL
- **Tool management**: mise (Go 1.23, Node 22)
- **Containerization**: Docker + Docker Compose (primary development and deployment method)
- **CI/CD**: GitHub Actions
- **Source control**: GitHub (github.com/ivsanmendez/ControlDeGastos)

## Architecture
- Monorepo: Go API at root, React SPA in `web/`
- Container-first: all services run via Docker Compose
- Go API serves JSON endpoints, React SPA consumes them
- PostgreSQL as the persistent data store

## Project Goals
- Track personal/household expenses
- Categorize and visualize spending
- Simple, maintainable codebase