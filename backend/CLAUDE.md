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

Schema decisions were worked through as a guided design exercise. Key decisions and their rationale are documented below.

### Design Decisions

**Archive strategy**: `is_archived bool` flag on each model (Project, Area, Resource). No separate archives table.
- Rationale: keeps archive/unarchive trivially reversible (just a field update); archive views are per-category in the UI so there is no need for a unified cross-table query.

**Project-Area relationship**: Non-nullable FK — every Project must belong to an Area.
- Rationale: enforced at the application layer (Projects can only be created from within an Area's page). Users are guided to create an Area first on onboarding.

**Resources**: Standalone items (belong only to a user, not to any Area or Project). Have a `type` field (Article | Book | Video | Note) for differentiated UI treatment. Books have a nullable `isbn` field. `content` was renamed to `description` for consistency.

**Project-Resource relationship**: Many-to-many via a join table (`project_resources`). One Resource can be linked to many Projects and vice versa.
- Join table columns: `project_id`, `resource_id`, `linked_at` (timestamp for ordering)
- Implemented as a manual `ProjectResource` struct (no GORM `many2many` association tag) because `linked_at` is a custom column beyond what GORM's auto join table supports.

**Project prerequisites**: Self-referential nullable FK on Project. A Project can optionally depend on another Project completing first.
- Cross-project prerequisite cycles must be detected and rejected at the API layer (not enforced by the DB).

**ToDo prerequisites**: Self-referential nullable FK on ToDo. A ToDo can optionally depend on another ToDo in the same Project.
- Prerequisites must reference a ToDo within the same Project — enforced at the API layer.
- Cycle detection (A → B → A) is an API concern, not a database constraint.

**ToDo CompletedAt**: Set by the system when Status transitions to `Closed`. Cleared when Status transitions away from `Closed`. Not user-settable directly.

### Schema

```
users                (done — auth layer)
  id                 uuid PK
  username           text
  email              text unique
  password           text (bcrypt hashed)
  created_at         timestamptz

areas                (ongoing responsibilities, no deadline)
  id                 uuid PK
  user_id            uuid FK → users.id  NOT NULL
  name               text
  description        text
  is_archived        bool default false
  created_at         timestamptz

projects             (active work with an intended outcome)
  id                 uuid PK
  user_id            uuid FK → users.id  NOT NULL
  area_id            uuid FK → areas.id  NOT NULL
  prerequisite_id    uuid FK → projects.id  nullable (self-referential)
  name               text
  description        text
  is_archived        bool default false
  created_at         timestamptz

todos                (tasks within a project)
  id                 uuid PK
  project_id         uuid FK → projects.id  NOT NULL
  prerequisite_id    uuid FK → todos.id  nullable (self-referential, same project only)
  title              text
  description        text
  due_date           timestamptz nullable
  status             text (Inactive | Working | Blocked | Closed)
  completed_at       timestamptz nullable (set by system on Closed, cleared on reopen)
  created_at         timestamptz

resources            (reference material / topics of interest)
  id                 uuid PK
  user_id            uuid FK → users.id  NOT NULL
  title              text
  description        text
  isbn               text (nullable, books only)
  source_url         text (nullable)
  type               text (Article | Book | Video | Note)
  is_archived        bool default false
  created_at         timestamptz

project_resources    (join table — many-to-many between projects and resources)
  project_id         uuid FK → projects.id
  resource_id        uuid FK → resources.id
  linked_at          timestamptz
  PRIMARY KEY (project_id, resource_id)
```

## Local Dev Runbook

All commands run from the `backend/` directory.

**1. Copy and fill in environment config (first time only)**
```
cp app.env.example app.env
# edit app.env with your values
```

**2. Start Postgres**
```
docker compose up -d
```

**3. Run migrations**
```
go run ./migrate/migrate.go
```

**4. Seed the database**
```
go run ./fixtures/seed.go
```

**4a. Clear the database**
```
go run ./fixtures/seed.go -teardown
```

**5. Start the API server**
```
go run main.go
```

## Current State

- Auth endpoints implemented: `POST /api/auth/register`, `POST /api/auth/login`
- JWT middleware for protected routes is in place
- PostgreSQL connection + GORM auto-migrate wired up
- Docker Compose local dev environment ready
- **Complete**: GORM models — `area.go`, `project.go`, `todo.go`, `resource.go` (includes `ProjectResource` join table)
- **Complete**: `migrate/migrate.go` updated to register all PARA models
- **Complete**: `fixtures/seed.go` — loads from `fixtures/seed_data.json`, seeds user + areas + projects + todos with prerequisites and due dates. Verified working.
- **Next**: CRUD endpoints for all PARA resources

## Seed Script Design

- `fixtures/seed_data.json` holds all seed data as plain JSON — edit this to change seed content
- `fixtures/seed.go` reads the JSON, converts to internal seed structs, and inserts into the DB
- Date fields (`due_date`, `completed_at`) are stored as day offsets (`due_in_days`, `completed_in_days`) in JSON and computed at runtime relative to when the seed runs
- JSON structs (`todoJSON`, `projectJSON`, `areaJSON`) and their `toSeed()` methods live in `seed.go` — they are seed-specific and do not belong in the models package
- Idempotent: checks for the seed user by email before inserting anything
