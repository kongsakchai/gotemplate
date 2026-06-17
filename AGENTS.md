# AGENTS.md

## Project

A **template** for Go web service APIs using Echo v5 + sqlx. Designed to be cloned via `gonew`:

```sh
go install golang.org/x/tools/cmd/gonew@latest
gonew github.com/kongsakchai/gotemplate github.com/yourname/projectname
```

## Commands

| Command | What it does |
|---------|-------------|
| `make test` | `go test -v ./... \| ./.script/colorize` (colorized output) |
| `make testcover` | Same as test with `-cover` |
| `make coverage` | Generates `coverage.out` and opens HTML report |
| `make init` | Installs mockery v3.7.3 + go-swagger |
| `make gen-mock` | Runs `mockery` to generate mocks |
| `make gendocs` | `swagger generate spec -o ./docs/swagger.yaml --scan-models --tags=docs` |
| `make docs` | Serves swagger UI from `./docs/swagger.yaml` |

Run a single package test: `go test -v ./pkg/config/`

## Architecture

- `pkg/` — shared utilities (config, database, cache, logger, httpclient, validator, errs, clock)
- `app/` — application layer: response helpers, Echo setup, middleware, domain modules
- `app/{domain}/` — each business domain as a package (e.g. `app/member/`)
- Business logic must live in `app/`, never in `pkg/`
- Entrypoint: `main.go` (not `cmd/` — Dockerfile references `./cmd/api` but that's a template artifact; update after `gonew`)

## Conventions (non-obvious)

- **`app.NotFound()` returns HTTP 400, not 404** — 404 is reserved for missing routes only
- **Minimize pointers** — use `(T, bool, error)` for existence checks instead of `(*T, error)`
- **Wrap errors** with `errs.From(err)` or `errs.New(...)` to preserve stack traces
- **Log through `echo.Context.Logger()`** — never use standalone `slog`/`log` directly
- **Handler + Service combined** — process functions are unexported methods on the handler struct, not a separate layer
- **One file per endpoint** — named after the action: `get_user_handler.go`, `create_user_handler.go`
- **Mockery generates `mock_*_test.go`** files in-source (same package, `_test.go` suffix)

## Test quirks

- Uses `github.com/labstack/echo/v5/echotest` for mocking Echo context
- In-memory SQLite (`modernc.org/sqlite`) for DB tests without infrastructure
- `github.com/DATA-DOG/go-sqlmock` for mocking SQL queries
- Mocks live in `_test.go` files alongside source — add `//mockery:generate: true` above any interface that needs mocking

## Config

Env vars loaded by `caarlos0/env/v11`. Prefix-based override: `ENV=LOCAL` enables `LOCAL_DATABASE_URL`, etc. See `.env.example` for all vars. Migration SQL files go in `migrations/` named `{version}_{name}.up.sql` / `.down.sql`.

## Skills

Load `.skills/gotemplate-guideline/SKILL.md` when writing or editing handlers, tests, or domain modules.
