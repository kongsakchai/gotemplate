# 🚀 Go Template

🎉 This is a template for creating a Go web service API project. Intended to be used as a starting point for creating a new Go web service API project and be a guideline for the project structure.

## 🔥 Usage

- Install `gonew`

```sh
go install golang.org/x/tools/cmd/gonew@latest
```

- Create a new project

```sh
gonew github.com/kongsakchai/gotemplate github.com/yourname/projectname
```

---

## 🌱 Project structure

```sh
./
├── app
│   └── middleware
├── cache
├── config
├── database
├── errs
├── httpclient
├── logger
├── middleware
├── migrations
└── validator
```

- **app** Application layer and business logic.

- **cache** Cache connectors, such as Redis.

- **config** Application configuration files and environment variable management.

- **database** Database connectors and setup, e.g., MySQL or PostgreSQL.

- **errs** Custom error types and centralized error handling for error tracking.

- **httpclient** HTTP client utilities for calling external services or APIs.

- **logger** Logging configuration and shared logger instances.

- **middleware** HTTP middleware for request processing, such as authentication, authorization, and logging.

- **validator** Request data validation logic, e.g., using [go-playground/validator](https://github.com/go-playground/validator).

- **migrations** Database migration files (.sql) for schema changes, e.g., using [kongsakchai/simple-migrate](https://github.com/kongsakchai/simple-migrate).

---

## 📚 Guideline Template

### Package `app/`

This is the main package we will focus on. Business logic and application layers should reside in this package, with each module clearly separated. Example:

```sh
./
└── app
    ├── user
    └── admin
        ├── admin.go
        └── handler.go
```

- app/register
- app/booking
- app/product

> [!CAUTION]
> Business logic should not be written in any package other than `app/`.

#### API Response

`app/app.go` provides helpers for API responses:

```go
type Response struct {
	Code    string `json:"code"` // business code
	Success  string `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}
```

**OK 200**

```go
app.Ok(ctx echo.Context, data any, msg ...string) error
// usage
app.OK(ctx, "data","success")
```

```yaml
status: 200
body: { 'code': '0000', 'success': true, 'data': 'data', 'message': 'success' }
```

**Created 201**

```go
app.Created(ctx echo.Context, data any, msg ...string) error
// usage
app.Created(ctx, "data","success")
```

```yaml
status: 201
body: { 'code': '0000', 'success': true, 'data': 'data', 'message': 'success' }
```

**API Error Response — `app/error.go`**

```go
type Error struct {
	HTTPCode int    // HTTP status code: 500, 400, 401, 403, 409
	Code     string // Business code
	Message  string
	Data     any
	Err      error  // Used for server-side logging only
}
```

```go
app.InternalServer(code string, msg string, err error, data ...any) app.Error
app.BadRequest(code string, msg string, err error, data ...any) app.Error
app.NotFound(code string, msg string, err error, data ...any) app.Error
app.Unauthorized(code string, msg string, err error, data ...any) app.Error
app.Forbidden(code string, msg string, err error, data ...any) app.Error
app.Conflict(code string, msg string, err error, data ...any) app.Error
```

> [!NOTE]
>
> - Why I do not use HTTP `404` for data not found ? `404` is a standard HTTP error that represents a missing endpoint or resource. Using it for missing data can cause confusion between “data not found” and “route not found”, and it also adds unnecessary complexity on the client side.
> - Why I use HTTP `400` for data not found ? I see missing data as something that usually results from an invalid or incorrect request from the client, while the system itself is still operating normally.

**500 Internal Server Error**

```go
app.Fail(ctx echo.Context, err app.Error) error
// usage
app.Fail(ctx, app.InternalServer(app.ErrInternalCode, app.ErrInternalMsg, err))
```

```yaml
status: 500
body: { 'code': '9999', 'success': false, 'message': 'internal error' }
```

**Global Error Handler**

```go
app.ErrorHandler(err error, ctx echo.Context)
// Usage
echoApp.HTTPErrorHandler = app.ErrorHandler // already configured in route.go
```

You can return an `app.Error` directly from a handler:

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

### Package `/app/middleware`

A helper package for HTTP middleware used in request processing,
such as authentication, authorization, logging, and request tracing.

**app/middleware/refid.go**

Middleware for managing a **reference ID** to make log tracing easier.
The header key can be configured via an environment variable.

```env
HEADER_REF_ID_KEY=
```

If the reference ID is not present in the request header, a new one will be generated using [github.com/google/uuid](https://github.com/google/uuid).

**app/middleware/logger.go**

Middleware for logging API request and response data.

### Package `cache/`

A helper package for interacting with caching systems. It includes utilities such as a Redis client factory. You may also integrate other caching solutions, such as [github.com/patrickmn/go-cache](https://github.com/patrickmn/go-cache).

### Package `config/`

All configuration should be read and stored as structs within this package. You can differentiate environments using the `ENV` variable and per-environment prefixes:

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

### Package `database/`

A package for creating database connectors. Files should be organized by database type, for example:

- database/mysql.go
- database/postgres.go
- database/mongo.go

### Package `errs/`

A helper package for error handling and error tracing:

```go
newErr := errs.Wrap(/* normal error */ err)
// OR
newErr := errs.New("some error")

fmt.Println(newErr.Error())
```

```
error: msg at (file.go:line) package.function
```

### Package `/httpclient`

A helper package for interact with external API using HTTP client. Contain a function to call external API and return the ressult

```go
type Response[T any] struct {
	Code    int // http code
	Data    T
	RawData []byte // raw rasponse
}
```

```go
httpclient.Get[Resp any](ctx context.Context, client *Client, url string, headers ...http.Header) (Response[Resp], error)
httpclient.Post[Resp any](ctx context.Context, client *Client, url string, payload any, headers ...http.Header) (Response[Resp], error)
httpclient.Put[Resp any](ctx context.Context, client *Client, url string, payload any, headers ...http.Header) (Response[Resp], error)
httpclient.Delete[Resp any](ctx context.Context, client *Client, url string, payload any, headers ...http.Header) (Response[Resp], error)
```

### Package `/logger`

A helper package for configuring the application logger.
You can control the log level, format, and enable/disable logging via environment variables.

```env
LOG_ENABLE=true
LOG_HTTP_ENABLE=true
LOG_LEVEL=debug|info|warning|error|critical
LOG_FORMAT=text|json
```

Sensitive data masking (such as passwords, tokens, or PII) can be configured in `logger/replace.go`.

### Package `/validator`

A package for defining validation rules for requests or structs using validation tags,
powered by [go-playground/validator](https://github.com/go-playground/validator).

### Folder `/migrations`

A folder containing SQL files for database migrations or schema updates.
Migration files follow this naming convention:

```text
version_name.up.sql
version_name.down.sql
```

**Example:**

```text
0001_init_schema.up.sql
```

Migration behavior can be configured via environment variables:

```env
MIGRATION_ENABLE=true
MIGRATION_DIR=./migrations
MIGRATION_VERSION=0001
MIGRATION_TABLE_NAME=schema_migrations # table used to store migration logs
```

- If `MIGRATION_VERSION` is not specified, the latest version will be used
- If `MIGRATION_TABLE_NAME` or `MIGRATION_DIR` is not specified, default values will be applied

### Recommended patterns

When using this Go template, I recommend the following patterns:

- **Storage pattern / Repository pattern** for managing database or external API interactions to separate concerns and improve testability.
- **Combine Handler with Service** I don't see the necessity to separate Handler from Service, as it may overcomplicate the code, especially for small to medium projects. Combining them reduces file count and improves code clarity and maintainability. However, I recommend breaking down Handler into smaller functions for better organization:
  - `Handle` function: manages HTTP requests
  - `Process` function: handles business logic (Service layer)

Example:

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

- **One file per endpoint** for clarity and easier maintenance. In larger projects, organizing files by endpoint improves code organization and makes features easier to locate and modify.
- **Separate modules by business domain** for better organization and maintainability. Domain-driven module separation improves code clarity and reduces cognitive load.
- **File naming** should represent the responsibility and purpose of the file.
- **Error handling** Use centralized error handling by creating custom error types and leveraging the global error handler to manage all errors in one place. This keeps code clean and simplifies maintenance.

### Testing

This project uses [testify](https://github.com/stretchr/testify) for testing. The `app` package provides a helper for mocking Echo context.

**Mocking Echo Context — `github.com/labstack/echo/v5/echotest`**

```go
ctx := echotest.ContextConfig{
	Headers: http.Header{
		echo.HeaderContentType: []string{echo.MIMEApplicationJSON},
	},
	JSONBody: []byte(`{"firstName":"john","lastName":"doe"}`),
}.ToContext(t)

ctx, rec := echotest.ContextConfig{
	Headers: http.Header{
		echo.HeaderContentType: []string{echo.MIMEApplicationJSON},
	},
	JSONBody: []byte(`{"firstName":"john","lastName":"doe"}`),
}.ToContextRecorder(t)
```

The function returns:

- `echo.Context` — for passing to handlers
- `*httptest.ResponseRecorder` — for asserting HTTP response

**Example**

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

**Dependecy injection**

Use `mockery` for generating mocks of interfaces. This allows you to easily create mock implementations of your interfaces for testing.

- add directive `//mockery:generate: true` to interface:

```go
//mockery:generate: true
type Storager interface {
    Users() ([]User, error)
    UserByName(name string) (User, error)
    CreateUser(user User) error
}
```

- Install mockery:

```bash

go install github.com/vektra/mockery/v2@latest
```

- Run mock generation: `mockery`

**Database testing**

- Use `modernc.org/sqlite` for testing database interactions. This allows you to create an in-memory SQLite database for testing purposes, which is fast and does not require any setup.

```go
import (
		_ "modernc.org/sqlite"
		"github.com/jmoiron/sqlx"
)

func TestDatabase(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	// Run migrations or setup schema here

	// Perform database operations and assertions
}
```

- Use `github.com/DATA-DOG/go-sqlmock` for testing database interactions without an actual database. This allows you to mock database queries and responses, making it easier to test your database logic in isolation.

```go
import (
		"github.com/DATA-DOG/go-sqlmock"
)

func TestDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Setup expected queries and responses here

	// Perform database operations and assertions
}
```

> https://github.com/DATA-DOG/go-sqlmock
