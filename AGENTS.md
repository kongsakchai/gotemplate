# AGENTS.md - Agent Coding Guidelines

## Project Overview

This is a Go backend project using Echo framework. It provides a template for building REST APIs with MySQL/SQLite database support, Redis caching, and comprehensive error handling.

> **Important**: Business logic should only be written in the `app/` package. Other packages (`config`, `database`, `cache`, `logger`, `validator`, etc.) are for infrastructure concerns only.

## Build, Lint, and Test Commands

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./app/...

# Run a single test file
go test -v ./app/app_test.go

# Run a single test function
go test -v -run TestAppResponse ./app/

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run integration tests (requires Docker/testcontainers)
go test -v -tags=integration ./...
```

### Build Commands

```bash
# Build the application
go build -o bin/server .

# Run the application
go run .
```

### Linting

This project uses standard Go tooling. Ensure code is formatted with:

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run all checks (fmt + vet + test)
go build ./... && go vet ./... && go test ./...
```

## Code Style Guidelines

### Imports

- Use Go's standard import organization:
  1. Standard library packages
  2. Third-party packages
  3. Internal packages

Example:
```go
import (
    "context"
    "fmt"
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"

    "github.com/kongsakchai/gotemplate/app"
)
```

- Blank import (`_`) for driver registration at the end of import block

### Formatting

- Use `gofmt` or an IDE with Go formatting support
- Maximum line length: use IDE wrapping or let `gofmt` handle it
- Leave a blank line between import groups and code

### Types

- Use explicit types for function parameters and return values
- Use `any` for generic data fields in structs (e.g., `Data any`)
- Use struct tags for JSON serialization: `json:"fieldName"`

Example:
```go
type Response struct {
    Code    string `json:"code"`
    Success bool   `json:"success"`
    Message string `json:"message,omitempty"`
    Data    any    `json:"data,omitempty"`
}
```

### Naming Conventions

- **Variables**: camelCase (e.g., `appPort`, `userName`)
- **Constants**: PascalCase for exported, camelCase for unexported (e.g., `SuccessCode`, `maxRetries`)
- **Functions**: PascalCase for exported, camelCase for unexported
- **Types/Structs**: PascalCase (e.g., `App`, `CreateUserRequest`)
- **Interfaces**: PascalCase, typically with `er` suffix for single-method interfaces (e.g., `Storager`)
- **JSON fields**: snake_case in tags, PascalCase in struct (e.g., `FirstName string \`json:"firstName"\``)
- **Package names**: short, lowercase, no underscores (e.g., `app`, `config`, `errs`)
- **File naming**: represent the responsibility and purpose (e.g., `user_handler.go`, `storage.go`)

## API Response Helpers

Use helpers from `app/app.go`:

```go
// OK 200
app.Ok(ctx, data, message...) 

// Created 201
app.Created(ctx, data, message...)

// Error response
app.Fail(ctx, app.Error)
```

Response structure:
```go
type Response struct {
    Code    string `json:"code"`    // business code
    Success bool   `json:"success"` 
    Message string `json:"message,omitempty"`
    Data    any    `json:"data,omitempty"`
}
```

## Error Handling

- Use custom `app.Error` type for API errors with HTTP status codes:
```go
type Error struct {
    HTTPCode int    // HTTP status code: 500, 400, 401, 403, 409
    Code     string // Business code
    Message  string
    Data     any
    Err      error  // Used for server-side logging only
}
```

- Error factory functions in `app` package:
  - `app.InternalError(code, msg, err, data...)` - 500
  - `app.BadRequest(code, msg, err, data...)` - 400
  - `app.NotFound(code, msg, err, data...)` - 400 (not 404 - see note below)
  - `app.Unauthorized(code, msg, err, data...)` - 401
  - `app.Forbidden(code, msg, err, data...)` - 403
  - `app.Conflict(code, msg, err, data...)` - 409

> **Note on 404 vs 400**: 404 is reserved for missing endpoints/resources. Using 400 for "data not found" avoids confusion between "route not found" and "data not found". Missing data is treated as an invalid client request.

- Check for empty errors using `err.IsEmpty()`
- Use global error handler: `app.ErrorHandler(err, ctx)` (already configured in `route.go`)

## Handler Pattern

Combine Handler with Service - no need to separate them. Break Handler into smaller functions:

```go
type handler struct {
    storage Storager
}

func NewHandler(storage Storager) *handler {
    return &handler{storage: storage}
}

func (h *handler) HandleRequest(ctx echo.Context) error {
    // 1. Bind request
    var req RequestStruct
    if err := ctx.Bind(&req); err != nil {
        return app.BadRequest("4001", "invalid request body", err)
    }

    // 2. Validate request
    if err := ctx.Validate(&req); err != nil {
        return app.BadRequest("4002", "validation error", err, err)
    }

    // 3. Process business logic (service layer)
    result, err := h.processRequest(req)
    if err != nil {
        return err // Return app.Error directly
    }

    // 4. Return response
    return app.Ok(ctx, result)
}

func (h *handler) processRequest(req RequestStruct) (*Response, app.Error) {
    // business logic here
}
```

### One file per endpoint

For clarity and easier maintenance, organize by endpoint:
```
app/user/
  user.go
  get_user_handler.go
  create_user_handler.go
  storage.go
```

### Storage/Repository Pattern

Use dependency injection for storage to improve testability:

```go
//mockery:generate: true
type Storager interface {
    Users() ([]User, error)
    UserByName(name string) (User, error)
    CreateUser(user User) error
}
```
Run mock generation: `go generate ./...`

## Configuration

- Use the `config` package with environment variables
- Use struct tags for env var mapping: `env:"VAR_NAME" envDefault:"default"`
- Load config with prefix: `config.Load("APP")`

Support environment-specific config:
```env
ENV=LOCAL|DEV|PROD
LOCAL_DATABASE_URL=
DEV_DATABASE_URL=
PROD_DATABASE_URL=
```

Example:
```go
type Database struct {
    URL string `env:"DATABASE_URL"`
}
```

## Validation

- Use `go-playground/validator` with `validate` tags
- Validate in handlers using `ctx.Validate(&struct)`

```go
type CreateUserRequest struct {
    FirstName string `json:"firstName" validate:"required"`
    LastName  string `json:"lastName" validate:"required"`
    Age       int    `json:"age" validate:"required,gte=0,lte=130"`
}
```

## Testing

- Use `github.com/stretchr/testify` (assert for checks, require for fatal failures)
- Name test files: `*_test.go`
- Use subtests with `t.Run()` for better organization:

```go
func TestExample(t *testing.T) {
    t.Run("should do something", func(t *testing.T) {
        // test code
    })
}
```

- Use `require.NoError` for setup/assertions that must pass
- Use `assert` for assertions that don't need to halt

## Database

- Use `sqlx` for database operations
- Follow repository pattern with interfaces
- Organize by database type:
  - `database/mysql.go`
  - `database/postgres.go`
  - `database/sqlite.go`

## Migrations

- Store in `/migrations` folder
- Naming convention:
  ```
  version_name.up.sql
  version_name.down.sql
  ```

Configuration via environment variables:
```env
MIGRATION_ENABLE=true
MIGRATION_DIR=./migrations
MIGRATION_VERSION=0001
MIGRATION_TABLE_NAME=schema_migrations
```

## Other Packages

### errs/ - Error Wrapping
```go
newErr := errs.Wrap(err)
// OR
newErr := errs.New("some error")
// Output: error: msg at (file.go:line) package.function
```

### httpclient/ - External API Calls
```go
type Response[T any] struct {
    Code    int
    Data    T
    RawData []byte
}

httpclient.Get[Resp any](ctx, client, url, headers...)
httpclient.Post[Resp any](ctx, client, url, payload, headers...)
```

### logger/ - Logging
```env
LOG_ENABLE=true
LOG_HTTP_ENABLE=true
LOG_LEVEL=debug|info|warning|error|critical
LOG_FORMAT=text|json
```
Sensitive data masking configured in `logger/replace.go`.

### cache/ - Redis caching
Use for Redis client factory and caching utilities.

### Middleware
- Place custom middleware in `app/middleware/` or `middleware/`
- Use Echo's `Use()` for global middleware
- Reference ID tracking via `X-Ref-ID` header (configurable via `HEADER_REF_ID_KEY`)

## Project Structure

```
./
├── app/
│   ├── apperror/          # Application errors
│   ├── example/           # Example handlers
│   └── middleware/        # Custom middleware
├── cache/                 # Redis caching
├── config/                # Configuration
├── database/              # Database connectors (mysql.go, postgres.go, etc.)
├── errs/                  # Error wrapping utilities
├── httpclient/            # HTTP client utilities
├── logger/                # Logging utilities
├── middleware/            # HTTP middleware
├── migrations/            # SQL migration files
├── validator/             # Validation utilities
└── main.go               # Entry point
```

### Module Organization

Separate modules by business domain:
```
app/
├── user/
│   ├── user.go
│   ├── get_user_handler.go
│   └── storage.go
├── admin/
│   ├── admin.go
│   └── handler.go
└── booking/
```
