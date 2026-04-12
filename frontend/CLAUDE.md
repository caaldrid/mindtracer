# frontend — Claude Code Context

## Stack

React 18, TypeScript, Vite, React Router, Axios

## Structure

```
src/
  pages/        Route-level components (MainPage, Register)
  components/   Shared UI components (Layout)
  data/         API client and data utilities
  App.tsx        Router setup
  main.tsx       Entry point
```

## API Client

`src/data/backend.ts` — `BackendAxiosCaller` class wraps Axios.
- Base URL from `VITE_APP_BACKEND_URL` env var
- Request interceptor attaches JWT from `localStorage` automatically
- Exposes `get`, `post`, `put`, `delete` methods

## Auth

- JWT stored in `localStorage` under key `accessToken`
- Backend expects `Authorization: Bearer <token>` header
- Interceptor in `backend.ts` handles this automatically

## Runbook

```bash
# Install deps
npm install

# Dev server
npm run dev

# Build
npm run build
```

Set `VITE_APP_BACKEND_URL` in a `.env` file pointing at the Go API (default: `http://localhost:8080/api`).

## Current State

- Routing set up (React Router)
- Register page — form + POST to `/auth/register`
- Main page — placeholder
- No login page yet
- No state management library yet

## Code Standards

This is a portfolio project. Code quality matters. Flag violations when you see them — don't silently go along with bad patterns.

**TypeScript**
- No `any` — use proper types or generics
- Define types/interfaces for all API response shapes
- Props must be typed — no untyped component props

**React**
- Functional components only — no class components
- Keep components focused — if a component is doing too much, flag it
- No business logic in JSX — extract to a function or hook
- Avoid inline object/array literals in JSX props where it causes unnecessary re-renders
- Custom hooks for reusable stateful logic — don't copy/paste `useState` + `useEffect` blocks

**State and side effects**
- `useEffect` dependencies must be complete and correct — flag missing deps
- Don't fetch data directly in components — keep API calls in `data/` or a custom hook
- No `any` casts to silence TypeScript errors

**Error handling**
- API errors must be handled visibly — don't silently swallow them in a `catch` with just a `console.error`
- User-facing errors should render something in the UI, not just log to console

**General**
- Named exports over default exports for everything except route-level pages
- Consistent file naming — components in PascalCase, utilities in camelCase

## Next Up

- Login page
- Protected routes (redirect to login if no token)
- PARA resource views (Projects, Areas, Resources, Archives)
