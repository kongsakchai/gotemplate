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

---

### Package `/middleware`

A helper package for HTTP middleware used in request processing,
such as authentication, authorization, logging, and request tracing.

**middleware/refid.go**

Middleware for managing a **reference ID** to make log tracing easier.
The header key can be configured via an environment variable.

```env
HEADER_REF_ID_KEY=
```

If the reference ID is not present in the request header, a new one will be generated using `uuid`.

**middleware/logger.go**

Middleware for logging API request and response data.

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
