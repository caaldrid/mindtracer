# TODO

## In Progress
- [ ] Add `//go:build integration` tag to integration test files and update Makefile targets accordingly
- [ ] Add storage interfaces and implementations for PARA resources (Area, Project, Resources, Todo)

## Up Next
- [ ] Register PARA routes and handlers in server.go
- [ ] Update frontend login form to send `email` instead of `username` (backend now uses `FindByEmail`)
- [ ] Show login error feedback in UI — currently only console.error on failed login
- [ ] Resolve `docs/notes/jwt-spa-workflows.md` — pick auth pattern, fix login state persistence and route protection
- [ ] Add frontend .env.example documenting VITE_APP_BACKEND_URL
- [ ] Resolve `docs/notes/llm-integration-go-backend.md` — define AI layer approach before implementing
- [ ] Resolve `docs/notes/openapi-first-look.md` — pick design-first vs code-first, then draft openapi.yaml
- [ ] Add `make generate` target to run oapi-codegen (blocked on openapi.yaml)

## Done
