# mindtracer — Claude Code Context

## Project

Full-stack personal information management app built on the PARA method (Projects, Areas, Resources, Archives).
End goal: AI-powered productivity tool for ADHD users.

This is a portfolio project. Carlos is the developer — he makes the architectural and design decisions.
Claude's role is to assist, answer questions, and execute tasks when asked. Not to lead.

## Stack

| Layer | Tech |
|-------|------|
| Backend | Go, Gin, GORM |
| Database | PostgreSQL (Docker) |
| Auth | JWT + bcrypt |
| Frontend | React, TypeScript, Vite |
| Config | Viper |

## Repo Structure

```
backend/    Go REST API
frontend/   React + TypeScript SPA
```

## Architecture Decisions

ADRs live in `docs/designs/decisions/`. These are standing decisions — code must conform to accepted ADRs. If an ADR conflicts with a task, flag it before implementing.

## Current Status

See `TODO.md` in the repo root.

## Portfolio Guardrails

This project needs to read as Carlos's work, not AI-generated code. These rules are standing instructions — apply them without being asked.

**Never do unprompted:**
- Add comments that describe *what* code does — only add comments when the *why* isn't obvious from the code
- Create handlers, models, storage methods, or routes for features not explicitly requested
- Add config fields, env vars, or flags speculatively
- Write commit messages — Carlos writes his own
- Add abstractions, interfaces, or wrappers beyond what the current problem requires

**Always do:**
- Flag over-engineering before implementing — if a solution feels like more complexity than the problem needs, say so first
- Flag when a decision is being made that Carlos should own — don't just pick an approach and implement it
- Keep implementations minimal and scoped to what was asked

## Permissions

Claude may run the following without asking:
- Docker (start/stop containers)
- Build commands
- Test runner
- Migrations
- Package management

See `backend/CLAUDE.md` and `frontend/CLAUDE.md` for stack-specific runbooks.
