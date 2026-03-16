# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run all tests
go test ./...

# Run tests in a specific package
go test ./internal/auth/...
go test ./internal/handlers/...

# Run a single test by name
go test ./internal/auth/ -run TestValidJWT

# Build
go build ./cmd/chirpy/

# Run the server (requires .env with DB_URL and JWT_SECRET)
go run ./cmd/chirpy/
```

### Database (sqlc + goose)

SQL queries live in `sql/queries/`, schemas in `sql/schema/`. After modifying either, regenerate the Go code:

```bash
sqlc generate
```

Migrations use numbered files (`001_`, `002_`, ...) in `sql/schema/`. Run them with goose or apply manually via psql.

## Architecture

This is a Boot.dev course project — a Twitter-like HTTP API in Go with PostgreSQL.

**Request flow:** `main.go` → `app.Server` (wires deps) → `app.Routes()` (registers handlers) → `handlers.Handlers` (handles requests) → `database.Queries` (sqlc-generated DB layer)

**Key packages:**

- `internal/auth` — password hashing (argon2id), JWT creation/validation (`MakeJWT`, `ValidateJWT`), and bearer token extraction (`GetBearerToken`). Tokens use `chirpy-access` as the issuer and embed the user's UUID in both `Subject` and a custom `UserID` claim.
- `internal/handlers` — one file per resource (`users.go`, `chirps.go`, `authentication.go`). All handlers are methods on `Handlers`, which holds `*database.Queries`, `*metrics.PageStats`, and `jwtSecret`.
- `internal/database` — sqlc-generated. Never edit directly; regenerate from `sql/` instead.
- `internal/httpjson` — thin helpers: `httpjson.Respond(w, status, v)` and `httpjson.Error(w, status, msg)`.

**Environment variables required:**

- `DB_URL` — PostgreSQL connection string
- `JWT_SECRET` — secret key for signing JWTs

**Auth state:** Fully implemented. `LogIn` issues a JWT (max 1 hour, configurable via `expires_in_seconds` in the request body). `POST /api/chirps` requires a valid `Authorization: Bearer <token>` header — the user ID is derived from the token, not the request body.

**Chirp rules:** max 140 runes; profane words (`kerfuffle`, `sharbert`, `fornax`) are replaced with `****` by `CleanRequestBody`.

# Project Context

This project is for learning purposes.

# Learning Mode Guidelines

- Prioritize explaining the "why" behind decisions.
- If I make a mistake, guide me toward finding it myself before pointing it out.
- Encourage best practices for software engineering.
- Use interactive, step-by-step guidance.
