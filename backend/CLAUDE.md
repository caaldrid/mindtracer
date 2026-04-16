# backend — Claude Code Context

## Package Structure

```
main.go           Entry point — wires config, DB, storage, server
handlers/         HTTP layer (Gin). One file per resource.
storage/          DB layer (GORM). Storage interface + implementations.
models/           GORM model structs
setup/            Config loading (Viper), DB connect, migrations, seed
```

## Auth Flow

- `POST /api/auth/register` — create user, bcrypt password
- `POST /api/auth/login` — verify password, return JWT
- All other routes under `/api/` require `Authorization: Bearer <token>`
- JWT secret + lifespan come from config (Viper / `app.env`)

## Testing

Two test suites:

| Package | Type | How it runs |
|---------|------|-------------|
| `handlers/` | Unit | Mock storage, no DB needed |
| `storage/` | Integration | Real Postgres via testcontainers |

`handlers/run_test.go` and `storage/run_test.go` set up the test harnesses.
Integration tests spin up a Postgres container automatically — no manual DB setup needed.

## Runbook

```bash
# Start DB
docker compose -f docker-compose.yml up -d

# Run server
go run main.go

# Run server with seed data
go run main.go --seed

# Run all tests (integration tests require Docker)
go test ./...

# Run only unit tests (no Docker)
go test ./handlers/...
```

## Config

Copy `app.env.example` to `app.env`. Viper reads it automatically.

## Code Standards

This is a portfolio project. Code quality matters. Flag violations when you see them — don't silently go along with bad patterns.

**Error handling**
- Always return errors explicitly — no swallowing, no ignoring
- Wrap errors with context: `fmt.Errorf("createUser: %w", err)`
- Use sentinel errors for known failure cases (e.g. `ErrUserAlreadyExists`)
- Never panic in handler or storage code

**Control flow**
- Prefer early returns over nested `if` blocks
- Happy path should be the least-indented path

**Gin handlers**
- Use `ShouldBindJSON` not `BindJSON` — `ShouldBindJSON` returns an error, `BindJSON` calls `c.Abort()` internally which is harder to reason about
- Every handler must end with a `ctx.JSON` call — no code paths that return silently

**Interfaces**
- Define interfaces where they are consumed, not where they are implemented
- Keep interfaces small — only the methods the consumer actually needs

**General Go**
- No `any` unless unavoidable — use concrete types or generics
- No global state — wire dependencies through function arguments
- Exported names get a doc comment; unexported ones only if the logic isn't obvious

## Next Up

- PARA CRUD endpoints (Projects, Areas, Resources, Archives)
- Storage interface methods for each PARA resource
- Handler files for each PARA resource
