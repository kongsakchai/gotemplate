# 🚀 Go Template

A starting point for building Go web service APIs. Use it to scaffold a new project or follow its conventions as a project guideline.

## 🔥 Quick Start

Create a new project from this template:

```sh
go install golang.org/x/tools/cmd/gonew@latest

gonew github.com/kongsakchai/gotemplate github.com/yourname/projectname
cd projectname
make init   # installs mockery v3, go-swagger, sets up script permissions
```

---

## 🌱 Project Structure

```
./                          ← main package
├── internal/               ← business logic per domain
│   ├── auth/               ← User domain: models, errors, service, storager
│   └── todo/               ← Todo domain: models, errors, service, storager
├── app/                    ← shared application infra + per-module handlers
│   ├── authapp/            ← auth HTTP routes (/api/v1/auth/register, /login)
│   ├── todoapp/            ← todos HTTP routes (/api/v1/todos — protected by JWT)
│   ├── echo.go             ← Echo setup, global middleware stack
│   ├── echo_route.go       ← Route helpers: app.GET/POST/PUT/DELETE
│   ├── app.go              ← Response envelope {code, success, message?, data?}
│   ├── error.go            ← app.Error constructors (InternalError, BadRequest, etc.)
│   ├── const.go            ← Global business codes (0000, 1000, 9999…)
│   ├── middleware_*.go     ← Global middleware: Recovery, CORS, Tag, RefID, Logger
│   ├── middleware_jwt.go   ← AuthMiddleware (JWT verification → claims in context)
│   ├── middleware_error.go ← ErrorMiddleware (handler → app.Error mapping)
│   └── docs.go             ← Shared swagger:response definitions
├── docs/                   ← go-swagger output (docs/swagger.yaml)
├── migrations/             ← SQL migrations: 0001_init.up.sql
├── pkg/                    ← reusable infrastructure (shared across modules)
│   ├── cache/              ← Redis client wrapper
│   ├── clock/              ← time.Time abstraction for testing
│   ├── config/             ← env-based config loading (see §Config)
│   ├── database/           ← DB connectors (mysqldb, postgresdb, sqlitedb, sqldb)
│   ├── hash/               ← Password hashing / comparison
│   ├── httpclient/         ← Typed external HTTP calls: Get[T], Post[T]…
│   ├── jwttoken/           ← JWT Signer / Verifier
│   ├── logger/             ← slog setup (JSON/text), ErrorAttrs
│   ├── migrate/            ← golang-migrate runner
│   ├── null/               ← Nullable types helpers
│   ├── serror/             ← Error tracing + coded business errors
│   └── validator/          ← go-playground/validator setup
└── .script/                ← Generator tools
    ├── colorize            ← Test output colorizer
    └── example.http        ← HTTP examples
```

### High-level flow

```
HTTP handler (app/{domain}app)
    └─► Service (internal/{domain})
        └─► Storager adapter (adapter_storage.go)
            └─► pkg/database/*
```

Business code lives in `internal/{domain}`; HTTP layer lives in `app/{domain}app`; shared infrastructure lives in `pkg/`. A module can optionally grow a `consumer/{domain}consumer` for background work — both depend on `internal/{domain}`, never on each other.

---

## 💡 Conventions

Read the [gotemplate-guideline](.agents/skills/gotemplate-guideline/SKILL.md) skill for the full set of conventions. Key rules:

| Rule | Detail |
|------|--------|
| **Private concrete types** | Service, storage, and app structs are unexported. Constructors return the private pointer directly (never an interface). Export only contracts — interfaces + deps structs. |
| **Return bool, not nil pointers** | If a function signals existence, return `(Value, bool, error)` instead of `(*Value, error)`. |
| **Wrap external errors** | Always wrap 3rd-party errors with `serror.From(err)` so they carry trace metadata. |
| **Error mapping at the boundary** | Map business codes to `app.Error` inside each module's `handleError`. Register it once via `app.ErrorMiddleware(handleError)` on the group — never per handler. |
| **Log through context** | Use `ctx.Logger()` which already carries traceID and route `tag`. |

For more detail see the skill references: **Architecture**, **Private Types**, **Errors**, **App/Handler**, **Middleware**.

---

## 📦 pkg / Shared Infrastructure

All reusable, non-domain code belongs in `pkg/`:

| Package | Purpose |
|---------|---------|
| `config` | Env-based config loading. Struct fields use `env:"VAR"` + `envDefault:"…"`. Load with `config.Load(config.Env)`. |
| `logger` | slog setup (JSON/text). Sensitive data masking is configurable. |
| `serror` | Error tracing + coded business errors. Produces output like `error: <msg>, code: <code>, at: (file.go:line)`. |
| `validator` | go-playground/validator setup with JSON-field-aware error messages. |
| `jwttoken` | JWT `Signer` / `Verifier` implementation. |
| `hash` | Password hashing and comparison. |
| `httpclient` | Generic typed HTTP wrappers (`Get[T]`, `Post[T]`, `Put[T]`, `Delete[T]`). |
| `database` | Database connectors: `mysqldb`, `postgresdb`, `sqlitedb`, `sqldb` (base). |
| `cache` | Redis client wrapper. Also usable with any in-memory cache. |
| `clock` | Time abstraction — injectable timer for deterministic testing. |
| `null` | Nullable type utilities. |
| `migrate` | Runs `golang-migrate` against your database. |

---

## 🏗️ Domain Architecture

### Inside `internal/{domain}/`

Each business domain gets its own package. Files follow these naming rules:

| Prefix | Meaning | Example |
|--------|---------|---------|
| `{domain}.go` | Domain model, exported interfaces, error codes | `auth.go`, `todo.go` |
| `service.go` | Business logic (≤4 methods). Split into `service_{core}.go` when larger. | `service.go` |
| `adapter_{target}.go` | External dependency connector | `adapter_storage.go`, `adapter_cache.go` |
| `service_helper.go` | Cross-cutting utility used by multiple files | `service_helper.go` |

**Interface lifecycle:**
1. Define the contract in `{domain}.go` (e.g. `Storager`, `Servicer`).
2. Mark interfaces for mocking: `//mockery:generate: true`.
3. Implement the concrete struct as **private** (`type service struct`, `type storage struct`).
4. Constructor returns the private type: `func NewService(deps ServiceDeps) *service`.
5. Alias external types locally for mocking: `//mockery:generate: true type Hasher = hash.Hasher`.

See the [Private Types](.agents/skills/gotemplate-guideline/references/private-types.md) reference for full examples.

### Inside `app/{domain}app/`

The HTTP handler exposes one constructor and one route-register method:

```go
func (a *todoApp) RegisterRoute(echo *app.EchoApp) {
    g := echo.Group("/api/v1/todos",
        app.AuthMiddleware(a.verifier),
        app.ErrorMiddleware(handleError),
    )
    app.GET(g, "get-todos", "", a.getTodos)
    app.POST(g, "create-todo", "", a.createTodo)
    // …
}
```

- Route helpers live in `app/echo_route.go`: `app.GET`, `app.POST`, `app.PUT`, `app.DELETE`.
- Signature: `app.METHOD(router, "route_name", path, handler, middlewares…)`. Route names become the `tag` in logs and metrics.
- Bind + validate with `app.Request[Req](c)` — single call for binding and validation.

See [App/Handler](.agents/skills/gotemplate-guideline/references/app-handler.md) for details.

---

## 🔐 Authentication

`app.AuthMiddleware(verifier)` verifies a Bearer JWT and puts every claim into the Echo context.

| Header missing | Expired token | Other failure |
|---------------|---------------|---------------|
| `Unauthorized(10002, "missing token")` | `Unauthorized(10004, "token expired")` | `Unauthorized(10003, "unauthorized")` / `Unauthorized(10005, "invalid token")` |

Access the current user ID with `sub := c.Get("sub").(string)`.

---

## ⚠️ Error Handling

Errors flow in two stages:

### Stage 1 — `serror` (business layer)

In `internal/{domain}`, business errors are coded strings defined with `serror.NewCoded`:

```go
var ErrUserNotFound = serror.NewCoded("2001", "user not found")

// Return the error with trace metadata:
return ErrUserNotFound.Err()

// Wrap an external error with a traceable wrapper:
return serror.From(err)

// Attach data surfaced in the HTTP response:
ErrNotFound.Err().WithData(user.ID)
```

Serror produces log output like:
```
error: user not found, code: 2001, at: (adapter_storage.go:44) auth.(*storage).FindUserByUsername
```

### Stage 2 — `app.Error` (HTTP boundary)

Every `app/{domain}app` defines a `handleError` that maps business codes → `app.Error`:

```go
func handleError(err error) app.Error {
    if e, ok := serror.As(err); ok {
        switch e.Code() {
        case todo.ErrUserNotFound.Code, todo.ErrTodoNotFound.Code:
            return app.NotFound(e.Code(), e.Msg(), e, e.Data...)
        }
    }
    return app.InternalError(app.InternalErrorCode, app.InternalErrorMsg, err)
}
```

Register once per route group: `app.ErrorMiddleware(handleError)`. Handlers simply `return err`.

#### Global Business Codes

| Code | Meaning |
|------|---------|
| `0000` | Success |
| `1000` | Bad request (bind failure) |
| `1001` | Invalid request (validation failure) |
| `10002` | Missing Authorization header |
| `10003` | Unauthorized |
| `10004` | Token expired |
| `10005` | Invalid token |
| `9998` | Database not ready |
| `9999` | Internal server error |

Module-specific codes live in `internal/{domain}` (e.g. `2000`–`2002` for auth, `3000`–`3001` for todo).

#### Why NotFound → HTTP 400

Missing data is treated as invalid input rather than a missing resource. Using 400 prevents callers from confusing "data not found" with "route not found."

See [Errors](.agents/skills/gotemplate-guideline/references/errors.md) for the complete reference.

---

## 🧭 Middleware Stack

Global middleware runs in this order (defined in `app/echo.go`):

| Order | Middleware | Purpose |
|-------|-----------|---------|
| 1 | `Recover()` | Panic recovery → HTTP 500 |
| 2 | `CORS("*")` | Cross-origin headers |
| 3 | `TagMiddleware()` | Sets `tag` from the route name (used in logs/metrics) |
| 4 | `RefIDMiddleware()` | Reads `X-Ref-ID` header; generates UUID if absent. Stored as `traceID` |
| 5 | `LoggerMiddleware()` | Logs request + response with traceID/tag |

Inside handlers use `ctx.Logger()` — it already carries traceID and tag.

Per-route middleware (like `AuthMiddleware`, `ErrorMiddleware`) is attached during `RegisterRoute`.

See [Middleware](.agents/skills/gotemplate-guideline/references/middleware.md) for details.

---

## 📤 API Responses

Success helpers (`app/app.go`):

| Function | Status | Envelope |
|----------|--------|----------|
| `app.Ok(ctx, data, msg?)` | 200 | `{"code":"0000","success":true,"message":msg,"data":data}` |
| `app.Created(ctx, data, msg?)` | 201 | `{"code":"0000","success":true,"message":msg,"data":data}` |

Error helpers return `app.Error`; the middleware serializes them as:
`{"code":"<code>","success":false,"message":"<msg>","data":<optional>}`

Request validation failures emit `{"code":"1001","success":false,"message":"invalid request"}`.

---

## 📝 Request Validation

Define request structs with `json` + `validate` tags:

```go
type createTodoRequest struct {
    Name        string `json:"name" validate:"required"`
    Description string `json:"description"`
    Status      string `json:"status"`
}
```

Bind + validate in one call:

```go
req, err := app.Request[createTodoRequest](c)
```

Path params use the `param:"id"` tag and are bound automatically.

---

## 📄 API Documentation (Swagger)

API docs use go-swagger with `//go:build docs` build tags:

1. Annotate endpoints in `app/{domain}app/docs.go` with `swagger:route`, `swagger:parameters`, `swagger:response`.
2. Shared responses live in `app/docs.go`; meta declaration is in `docs.go` (root package main).
3. Run `make gendocs` → writes `docs/swagger.yaml`.
4. Run `make docs` → serves Swagger UI in the browser.

---

## 🧪 Testing

Stack: **testify** (assertions), **mockery v3** (interface mocks), **go-sqlmock** (DB tests), **echotest** (handler tests). No `testify/suite` — prefer plain `t.Run` subtests.

### Generating Mocks

Mark an interface with a comment:

```go
//mockery:generate: true
type Storager interface {
    FindUser(ctx context.Context, id string) (User, bool, error)
}
```

Generate: `make genmock` (writes `mock_test.go` next to the interface).

To mock an interface from another package, alias and re-mark:

```go
//mockery:generate: true
type Servicer = auth.Servicer
```

Generated mocks come with helpers like `newMockStorager(t)` pre-registered with `*testing.T`.

### Service Tests

Build a service with mocked dependencies:

```go
mock := newMockStorager(t)
sv := auth.NewService(auth.ServiceDeps{Storager: mock})
```

Verify business errors with `serror.As`:

```go
serr, ok := serror.As(err)
assert.True(t, ok)
assert.Equal(t, serr.Code(), auth.ErrUserNotFound.Code)
```

### Storage Tests

Use go-sqlmock:

```go
db, mock, _ := sqlmock.New()
sqlxDB := sqlx.NewDb(db, "sqlmock")
mock.ExpectQuery("SELECT.*FROM user WHERE username = ?").WithArgs("admin").
    WillReturnRows(sqlmock.NewRows([]string{"id","username"}).AddRow("u1","admin"))
```

### Handler Tests

Use echotest with `ContextConfig`:

```go
ctx, rec := echotest.ContextConfig{
    Headers: http.Header{echo.HeaderContentType: []string{echo.MIMEApplicationJSON}},
    PathValues: map[string]string{"id": "u1"},
}.ToContextRecorder(t)

ctx.Set("sub", "user-1")  // replicate middleware injection

handler := &todoApp{sv: mock}
err := handler.getTodos(&ctx)

require.NoError(t, err)
assert.Equal(t, http.StatusOK, rec.Code)

var resp app.Response
json.Unmarshal(rec.Body.Bytes(), &resp)
assert.True(t, resp.Success)
```

Run: `make test` (colorized). See [Testing](.agents/skills/gotemplate-guideline/references/testing.md) for full conventions.

---

## ⚙️ Configuration

Config lives in `pkg/config`. Struct fields use `env:"VAR"` and `envDefault:"value"` tags:

```go
type Config struct {
    App    AppConfig    `envPrefix:"APP_"`
    Database DatabaseConfig `envPrefix:"DATABASE_"`
    Log    LogConfig    `envPrefix:"LOG_"`
}
```

Env vars are parsed twice — unprefixed first, then with the `{ENV}_` prefix:

| Precedence | Example | Source |
|------------|---------|--------|
| Highest | `DEV_DATABASE_URL=…` | Prefixed by ENV |
| Middle | `DATABASE_URL=…` | Unprefixed |
| Lowest | `envDefault:"sqlite:///dev.db"` | Default value |

Load: `cfg := config.Load(config.Env)` where `config.Env` holds the current `ENV` value.

---

## 🔄 Migrations

SQL migration files live in `migrations/` following the `golang-migrate` convention:

```
migrations/
├── 0001_init_schema.up.sql
├── 0001_init_schema.down.sql
└── 0002_add_status.up.sql
```

Environment variables control behavior:

| Variable | Purpose |
|----------|---------|
| `MIGRATION_ENABLE` | `"true"` / `"false"` — enable/disable migrations at startup |
| `MIGRATION_SRC` | Path to migration files directory |
| `MIGRATION_DATABASE_URL` | Target database URL |
| `MIGRATION_VERSION` | Specific version to run |

---

## 🛠️ Makefile Commands

| Command | Purpose |
|---------|---------|
| `make init` | Install tooling (mockery v3, go-swagger), set script permissions |
| `make test` | Run all tests with colorized output |
| `make testcover` | Run tests showing coverage inline |
| `make coverage` | Write `coverage.out` + open HTML report in browser |
| `make genmock` | Regenerate all mocks (mockery v3) |
| `make gendocs` | Generate `docs/swagger.yaml` from annotations |
| `make docs` | Serve Swagger UI (`localhost:8080/docs` by default) |

---

*For deep-dive conventions, read the [gotemplate-guideline](.agents/skills/gotemplate-guideline/SKILL.md) skill.*
