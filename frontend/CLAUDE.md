# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this directory.
It lives at `frontend/CLAUDE.md` and provides React/TypeScript-specific context for Phase 3.
Shared context (about the developer, your role, the learning path) lives in the root `CLAUDE.md`.

---

# mindtracer Frontend — Phase 3 Context

## What This Project Is

mindtracer's frontend is a React/TypeScript web app that consumes the Go backend API built
in Phase 2. It is the Phase 3 portfolio project and the first full-stack integration point
on the learning path.

The UI provides a PARA-method interface — helping users capture and organize their Projects,
Areas, Resources, and Archives through a clean, intentional interface designed with ADHD
users in mind.

## Tech Stack

> Verify these details against the actual repo when starting Phase 3.
> Update this section to reflect what's actually in `package.json` and config files.

- **Language**: TypeScript
- **Framework**: React
- **Build tool**: (check `package.json` — likely Vite or Create React App)
- **Styling**: (check repo)
- **HTTP client**: (check repo — likely fetch or axios)
- **Testing**: (check repo)

## Architecture Principles

- **Stateless UI** — all persistent state lives in the backend. The frontend is a pure
  consumer of the API; nothing business-critical lives only in the browser.
- **Config via environment** — API base URL and any other env-specific values should come
  from environment variables, not be hardcoded.
- **Component boundaries** — keep UI components focused and composable. Shared logic
  belongs in hooks or utilities, not scattered across components.

## Current State

The frontend directory exists in the repo but has not been the focus of active development.
Phase 3 begins here — audit what's present, understand the existing structure, then build
out the PARA resource views against the API built in Phase 2.

**Starting checklist for Phase 3:**
- [ ] Audit existing frontend structure and dependencies
- [ ] Update this CLAUDE.md with the actual tech stack
- [ ] Confirm API base URL config and auth token handling strategy
- [ ] Plan UI flows for each PARA resource (Projects, Areas, Resources, Archives)
