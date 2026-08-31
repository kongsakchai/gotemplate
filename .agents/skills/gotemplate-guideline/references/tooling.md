# Tooling

Shared infra, config, migration, openapi docs, and make commands.

## pkg/ shared packages

Reusable, non-domain infrastructure lives in `pkg/`:

- `config` — env-based config loading
- `logger` — slog setup (JSON/text), `ErrorAttrs`
- `serror` — error tracing + business codes
- `validator` — request validation (go-playground), JSON-field-aware errors
- `jwttoken` — JWT `Signer` / `Verifier`
- `hash` — password hashing / compare
- `httpclient` — typed external HTTP calls (`httpclient.Get[Resp](ctx, client, url, ...)`)
- `database` — DB connectors (`mysqldb`, `postgresdb`, `sqlitedb`, `sqldb`)
- `cache` — Redis client
- `clock` — time abstraction
- `null` — nullable types
- `migrate` — migrations runner

## Config (pkg/config)

- Struct fields map to env vars with `env:"VAR"` + `envDefault:"..."` tags.
- Load with `config.Load(config.Env)` (note: `config.Env` is the current `ENV` value used as prefix; lowercase helpers like `env.Parse` are the library).
- Env vars are parsed twice — unprefixed first, then with the `{ENV}_` prefix — so precedence is `{ENV}_VAR` > `VAR` > `envDefault` (e.g. with `ENV=DEV`, `DEV_DATABASE_URL` wins over `DATABASE_URL`).
- `.env.example` documents every variable; keep it in sync when adding config.

```go
type App struct {
	Name    string `env:"APP_NAME" envDefault:"gotemplate"`
	Port    string `env:"APP_PORT" envDefault:"8080"`
	Version string `env:"APP_VERSION" envDefault:"0.0.1"`
}
```

## Migration

- Uses `golang-migrate`; SQL files live in `migrations/` as `NNNN_*.up.sql`.
- Controlled by config: `MIGRATION_ENABLE`, `MIGRATION_SRC`, `MIGRATION_DATABASE_URL`, `MIGRATION_VERSION`.

## Swagger / API docs

- go-swagger with `//go:build docs` build tag (`docs.go` declares `swagger:meta`; `app/docs.go` defines shared responses).
- Annotate routes/responses/params: `swagger:route`, `swagger:response`, `swagger:parameters`.
- Success responses wrap `app.SwaggerSuccessResponse` (`code: "0000"`, `success: true`).
- `make gendocs` writes `docs/swagger.yaml`; `make docs` serves the UI.

## Makefile

| Command | Purpose |
|---------|---------|
| `make init` | install tooling (mockery v3, go-swagger) |
| `make test` | run all tests (colorized) |
| `make testcover` | run tests with coverage |
| `make coverage` | write coverage profile to `coverage.out` and open the HTML report in the browser |
| `make genmock` | regenerate mocks with mockery |
| `make gendocs` | generate swagger spec |
| `make docs` | serve swagger UI |