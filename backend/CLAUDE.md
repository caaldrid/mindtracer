# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this directory.
It lives at `backend/CLAUDE.md` and provides Go/API-specific context for Phase 2.
Shared context (about the developer, your role, the learning path) lives in the root `CLAUDE.md`.

---

# mindtracer Backend — Phase 2 Context

## What This Project Is

mindtracer's backend is a PARA-method personal information management API written in Go.
It models the core data layer for The Toolbox and is designed to be consumed by the
React/TypeScript frontend in Phase 3, and an AI layer in Phase 4+.
It could later be exposed as an MCP server.

The PARA method organizes all information into four categories:
- **Projects** — active work with an intended outcome and a deadline
- **Areas** — ongoing responsibilities with no end date
- **Resources** — reference material and topics of interest
- **Archives** — inactive items from any of the above categories

User authentication (register, login, JWT) is already implemented. This phase focuses on
building out the PARA content endpoints and their data models.

## Tech Stack

- **Runtime/Language**: Go 1.24.3
- **Framework**: Gin (v1.10.1)
- **ORM**: GORM (v1.30.0)
- **Database**: PostgreSQL
- **Auth**: JWT (golang-jwt v5.2.2) + bcrypt
- **Config**: Viper (v1.20.1) via `app.env`
- **Local infra**: Docker Compose for PostgreSQL
- **Testing**: Go's built-in `testing` package + `net/http/httptest`

## Architecture Principles

This project is built with cloud native principles in mind from the start:

- **Stateless** — no session or user state stored in memory on the server. All persistent
  state lives in PostgreSQL. Any instance can handle any request.
- **Observable** — structured logging so there's something useful to debug in production.
- **Resilient** — proper error handling with meaningful HTTP status codes.
- **Loosely coupled** — components are modifiable without requiring changes elsewhere.
- **Config via environment** — no hardcoded connection strings, ports, or secrets. Use
  `app.env` locally, inject in production.

Go's standard patterns align naturally with these principles: explicit error returns keep
error handling visible and intentional; interfaces enable loose coupling without ceremony;
package-level separation enforces clear boundaries between concerns.

## Data Model

The exact schema will be worked through as a guided exercise — particularly
decisions like: Should archives be a separate table or a soft-delete flag? Should projects
link to areas? This is part of the learning, not something to prescribe upfront.

Starting point for discussion:

```
users           (done — auth layer)
  id            uuid PK
  username      text
  email         text unique
  password      text (bcrypt hashed)
  created_at    timestamptz

projects        (active work with an intended outcome)
  id            uuid PK
  user_id       uuid FK → users.id
  name          text
  description   text
  status        text  (active, completed, archived)
  created_at    timestamptz

areas           (ongoing responsibilities, no deadline)
  id            uuid PK
  user_id       uuid FK → users.id
  name          text
  description   text
  created_at    timestamptz

resources       (reference material / topics of interest)
  id            uuid PK
  user_id       uuid FK → users.id
  title         text
  content       text
  source_url    text (nullable)
  created_at    timestamptz

archives        (inactive items from any PARA category)
  id            uuid PK
  user_id       uuid FK → users.id
  original_type text  (project | area | resource)
  original_id   uuid
  archived_at   timestamptz
```

## Current State

- Auth endpoints implemented: `POST /api/auth/register`, `POST /api/auth/login`
- JWT middleware for protected routes is in place
- PostgreSQL connection + GORM auto-migrate wired up
- Docker Compose local dev environment ready
- **Not yet built**: GORM models, migrations, and CRUD endpoints for Projects, Areas,
  Resources, and Archives
