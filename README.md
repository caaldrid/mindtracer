# Mindtracer

A full-stack personal information management app built on the [PARA method](https://fortelabs.com/blog/para/) — a simple framework for organizing everything that matters into four categories: Projects, Areas, Resources, and Archives.

The end goal is an AI-powered productivity tool designed specifically for ADHD brains.

---

## What it does

- **Projects** — track active work with tasks, prerequisites, and linked resources
- **Areas** — maintain ongoing responsibilities with no defined end date
- **Resources** — collect reference material (articles, books, videos, notes) and link them to relevant projects
- **Archives** — soft-archive anything that's no longer active without losing it

---

## Stack

| Layer | Tech |
|-------|------|
| Backend | Go, Gin, GORM |
| Database | PostgreSQL |
| Auth | JWT + bcrypt |
| Frontend | React, TypeScript |
| Config | Viper |
| Local infra | Docker Compose |

---

## Status

Actively in development. Current phase: Go REST API — PARA data models and CRUD endpoints.

| Component | Status |
|-----------|--------|
| Auth (register, login, JWT) | Complete |
| PARA data models | Complete |
| CRUD endpoints | In progress |
| React frontend | Planned |
| AI layer | Planned |

---

## Local setup

> Documentation in progress. Full setup instructions will be added once the API layer is stable.

---

## Author

Carlos Aldridge — [github.com/caaldrid](https://github.com/caaldrid)
