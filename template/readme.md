# 🚀 Go Template

A starting point for Go web service API projects.

## 🔥 Usage

Install `gonew`:

```sh
go install golang.org/x/tools/cmd/gonew@latest
```

Create a new project:

```sh
gonew github.com/kongsakchai/gotemplate/template github.com/yourname/projectname
```

---

## 🌱 Project Structure

**Common packages**

```
./
├── cache
├── database
├── errs
├── httpclient
├── logger
├── pkg
└── validator
```

| Package      | Description                                                                                          |
| ------------ | ---------------------------------------------------------------------------------------------------- |
| `cache`      | Cache connectors (e.g., Redis)                                                                       |
| `database`   | Database connectors and setup (e.g., MySQL, PostgreSQL)                                              |
| `errs`       | Custom error types and centralized error handling                                                    |
| `httpclient` | HTTP client utilities for calling external services                                                  |
| `logger`     | Logging configuration and shared logger instances                                                    |
| `pkg`        | Small helper packages used across the project                                                        |
| `validator`  | Request validation logic using [go-playground/validator](https://github.com/go-playground/validator) |

**Template packages**

```
./
├── .script
├── app
│   ├── apperror
│   └── middleware
├── config
├── docs
└── migrations
```

| Package          | Description                                                                                         |
| ---------------- | --------------------------------------------------------------------------------------------------- |
| `app`            | Application layer and business logic                                                                |
| `app/apperror`   | Global error handler                                                                                |
| `app/middleware` | HTTP middleware (auth, logging, etc.)                                                               |
| `config`         | Configuration and environment variable management                                                   |
| `docs`           | API documentation (e.g., [go-swagger](https://github.com/go-swagger/go-swagger))                    |
| `migrations`     | SQL migration files (e.g., [simple-sql-migrate](https://github.com/kongsakchai/simple-sql-migrate)) |

---

## 📚 Package Guidelines

### `cache/`

Utilities for caching systems such as Redis. Other options like [go-cache](https://github.com/patrickmn/go-cache) can also be integrated.

### `database/`

Organize database connector files by type:

```
database/mysql.go
database/postgres.go
database/mongo.go
```

### `errs/`

Error handling and tracing utilities:

```go
err := errs.Wrap(err)
// or
err := errs.New("some error")

fmt.Println(err.Error())
// error: msg at (file.go:line) package.function
```

### `httpclient/`

Helpers for calling external APIs:

```go
type Response[T any] struct {
    Code    int
    Data    T
    RawData []byte
}

httpclient.Get[T](ctx, client, url, headers...)
httpclient.Post[T](ctx, client, url, payload, headers...)
httpclient.Put[T](ctx, client, url, payload, headers...)
httpclient.Delete[T](ctx, client, url, payload, headers...)
```

### `logger/`

Configure logging via environment variables:

```env
LOG_ENABLE=true
LOG_HTTP_ENABLE=true
LOG_LEVEL=debug|info|warning|error|critical
LOG_FORMAT=text|json
```

Sensitive data masking (passwords, tokens, PII) can be configured in `logger/replace.go`.

### `pkg/`

- `pkg/timer` — `Timer` interface for injecting a time source, useful for deterministic testing.
- `pkg/mockutil` — Test helpers and mocks for unit tests.

### `validator/`

Define validation rules using struct tags powered by [go-playground/validator](https://github.com/go-playground/validator).

---

## 📦 App Layer

### `app/`

All business logic lives here, organized by module:

```
app/
├── user/
└── admin/
    ├── admin.go
    └── handler.go
```

> [!CAUTION]
> Business logic must not be written outside the `app/` package.

### API Responses

**Success responses — `app/app.go`**

```go
type Response struct {
    Code    string `json:"code"`
    Success bool   `json:"success"`
    Message string `json:"message,omitempty"`
    Data    any    `json:"data,omitempty"`
}

app.Ok(ctx, data, msg...)      // 200
app.Created(ctx, data, msg...) // 201
```

```yaml
# Example
status: 200
body: { 'code': '0000', 'success': true, 'data': '...', 'message': 'success' }
```

**Error responses — `app/error.go`**

```go
type Error struct {
    HTTPCode int
    Code     string
    Message  string
    Data     any
    Err      error // server-side logging only
}

app.InternalServer(code, msg, err, data...)
app.BadRequest(code, msg, err, data...)
app.NotFound(code, msg, err, data...)
app.Unauthorized(code, msg, err, data...)
app.Forbidden(code, msg, err, data...)
app.Conflict(code, msg, err, data...)
```

> [!NOTE]
> This template uses `400` instead of `404` for missing data, since `404` should indicate a missing route, not a missing record. Missing data is typically the result of an invalid client request.

```go
// Usage
app.Fail(ctx, app.InternalServer(app.ErrInternalCode, app.ErrInternalMsg, err))
```

```yaml
status: 500
body: { 'code': '9999', 'success': false, 'message': 'internal error' }
```

### `app/apperror/`

Global error handler — already wired into the router:

```go
echoApp.HTTPErrorHandler = apperror.ErrorHandler
```

Handlers can return `app.Error` directly:

```go
func healthCheck(db *sqlx.DB) echo.HandlerFunc {
    return func(ctx echo.Context) error {
        if db.Ping() != nil {
            return app.InternalServer(app.ErrInternalCode, app.ErrDatabaseMsg, nil)
        }
        return app.Ok(ctx, nil, "healthy")
    }
}
```

### `app/middleware/`

**`middleware/refid.go`** — Injects a reference ID for log tracing. Reads from the request header or generates a new UUID if absent.

```env
HEADER_REF_ID_KEY=
```

**`middleware/logger.go`** — Logs API request and response data.

---

## ⚙️ Config

Read environment variables into typed structs. Use `ENV` to differentiate environments:

```env
ENV=LOCAL|DEV|PROD

LOCAL_DATABASE_URL=
DEV_DATABASE_URL=
PROD_DATABASE_URL=
```

```go
type Database struct {
    URL string `env:"DATABASE_URL"`
}
```

---

## 🗄️ Migrations

SQL migration files follow this naming convention:

```
version_name.up.sql
version_name.down.sql

# Example
0001_init_schema.up.sql
```

Configure via environment variables:

```env
MIGRATION_ENABLE=true
MIGRATION_DIR=./migrations
MIGRATION_VERSION=0001
MIGRATION_REPEAT=none
```

If `MIGRATION_VERSION` and `MIGRATION_REPEAT` are not set, the latest version is used. `MIGRATION_DIR` must not be empty.

---

## ✅ Recommended Patterns

**Repository pattern** — Separate database/external API interactions for better testability.

**Combined Handler + Service** — Avoid splitting Handler from Service unnecessarily. For small to medium projects, combining them reduces complexity. Break logic into two functions instead:

- `Handle` — manages the HTTP request
- `Process` — contains business logic

```go
type handler struct {
    storage Storager
}

func (h *handler) GetUserByID(ctx echo.Context) error {
    userID := ctx.Param("id")
    _, err := h.processGetUserByID(userID)
    if err != nil {
        return err
    }
    return nil
}

func (h *handler) processGetUserByID(userID string) (*User, error) {
    // business logic here
}
```

**One file per endpoint** — Improves readability and makes features easier to locate.

**Domain-separated modules** — Group code by business domain to reduce cognitive load.

**Descriptive file names** — File names should reflect their responsibility.

**Centralized error handling** — Use custom error types and the global error handler to keep code clean.

---

## 🧪 Testing

Uses [testify](https://github.com/stretchr/testify).

### Mocking Echo Context

```go
ctx, rec := echotest.ContextConfig{
    Headers: http.Header{
        echo.HeaderContentType: []string{echo.MIMEApplicationJSON},
    },
    JSONBody: []byte(`{"firstName":"john","lastName":"doe"}`),
}.ToContextRecorder(t)
```

Returns `echo.Context` for passing to handlers and `*httptest.ResponseRecorder` for asserting responses.

```go
func TestGetUser(t *testing.T) {
    ctx, rec := echotest.ContextConfig{
        Headers: http.Header{
            echo.HeaderContentType: []string{echo.MIMEApplicationJSON},
        },
        JSONBody: []byte(`{"firstName":"john","lastName":"doe"}`),
    }.ToContextRecorder(t)

    handler := NewHandler(mockStorage)
    err := handler.GetUser(ctx)

    require.NoError(t, err)
    assert.Equal(t, 200, rec.Code)
    assert.JSONEq(t, `{"code":"0000","success":true,"data":{...}}`, rec.Body.String())
}
```

### Dependency Injection with Mockery

Add the directive to the interface:

```go
//mockery:generate: true
type Storager interface {
    Users() ([]User, error)
    UserByName(name string) (User, error)
    CreateUser(user User) error
}
```

Install and run:

```sh
go install github.com/vektra/mockery/v2@latest
mockery
```

### Database Testing

**In-memory SQLite** (fast, no setup):

```go
import (
    _ "modernc.org/sqlite"
    "github.com/jmoiron/sqlx"
)

db, err := sqlx.Open("sqlite", ":memory:")
require.NoError(t, err)
defer db.Close()
```

**SQL mock** (no database required):

```go
import "github.com/DATA-DOG/go-sqlmock"

db, mock, err := sqlmock.New()
require.NoError(t, err)
defer db.Close()
```

> Reference: https://github.com/DATA-DOG/go-sqlmock
