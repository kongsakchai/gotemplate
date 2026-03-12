# AGENTS.md - Go Template Project Guidelines

## Overview
This is a Go web service API template using Echo framework. Business logic MUST reside in `app/` package only.

## Build, Lint, and Test Commands

### Running Tests
```bash
go test -v ./...                           # Run all tests
go test -cover ./...                       # With coverage
go test -coverprofile=coverage.out ./...   # Coverage report
go test -v ./path/to/package -run TestName # Single test
```

### Mock Generation
```bash
mockery  # Uses .mockery.yml config
```

## Code Style Guidelines

### Project Structure
```
./
├── app/              # Business logic only
├── cache/            # Redis/caching
├── config/           # Configuration
├── database/         # DB connectors
├── errs/             # Custom errors
├── httpclient/       # HTTP utilities
├── logger/           # Logging
├── middleware/       # HTTP middleware
├── migrations/      # SQL migrations
└── validator/       # Request validation
```

### Naming Conventions
- **Interfaces**: `Storager`, `Handler` (no "I" prefix)
- **Structs**: `User`, `CreateUserRequest` (PascalCase)
- **Functions**: `GetUser`, `createUser` (unexported = lowercase)
- **Variables**: `userID`, `config` (camelCase)
- **Database**: snake_case columns, PascalCase structs

### Imports
```go
import (
    "context"
    "net/http"

    "github.com/kongsakchai/gotemplate/app"
    "github.com/labstack/echo/v4"
)
```
- Standard library first, then external packages
- Use `goimports` for formatting

### Types & Structs
- Use `any` not `interface{}`
- Proper struct tags for JSON/validation:
```go
type CreateUserRequest struct {
    FirstName string `json:"firstName" validate:"required"`
    LastName  string `json:"lastName" validate:"required"`
    Age       int    `json:"age" validate:"gte=0,lte=130"`
}
```

### Error Handling
Use centralized `app.Error` type:
```go
// In handlers - return app.Error
return app.BadRequest("4001", "invalid request body", err)
return app.InternalError("5001", "failed to get user", err)
return app.Conflict("4003", "user already exists", nil)

// Global error handler
echoApp.HTTPErrorHandler = app.ErrorHandler
```

### Response Format
```go
app.Ok(ctx, data)           // 200 OK
app.Created(ctx, data)      // 201 Created
app.Fail(ctx, app.Error)   // Error response
```

### Handler Pattern
```go
type handler struct {
    storage Storager
}

// Public - handles HTTP
func (h *handler) CreateUser(ctx echo.Context) error {
    var req CreateUserRequest
    if err := ctx.Bind(&req); err != nil {
        return app.BadRequest("4001", "invalid request body", err)
    }
    return app.Ok(ctx, nil)
}

// Private - business logic, returns app.Error
func (h *handler) createUser(req CreateUserRequest) app.Error {
    err := h.storage.CreateUser(user)
    if err != nil {
        return app.InternalError("5002", "failed to create user", err)
    }
    return app.Error{} // Empty = success
}
```

### Testing
- Use `testify` (assert/require)
- Use `app.NewMockContext()` for Echo context mocking
- Use mockery with `//mockery:generate: true` directive
- Use `modernc.org/sqlite` for in-memory DB testing
```go
ctx, rec := app.NewMockContext(http.MethodPost, "/users", `{"firstName":"john"}`)
storage := newMockStorager(t)
storage.EXPECT().CreateUser(gomock.Any()).Return(nil)
```

### Configuration
```go
type Database struct {
    URL string `env:"DATABASE_URL"`
}
// Usage: config.Load("DEV") for DEV_DATABASE_URL
```

### Logging
```
LOG_ENABLE=true
LOG_HTTP_ENABLE=true
LOG_LEVEL=debug|info|warning|error|critical
LOG_FORMAT=text|json
```

### Database Migrations
- Place in `migrations/` directory
- Naming: `0001_init_schema.up.sql`, `0001_init_schema.down.sql`
- Use `github.com/kongsakchai/simple-migrate`

## Key Patterns

### Storage/Repository Pattern
```go
type Storager interface {
    UserByName(name string) (User, error)
    CreateUser(user User) error
}
```

### HTTP Client
```go
resp, err := httpclient.Get[ResponseType](ctx, client, url, headers)
```

## Dependencies
- **Echo v4**: Web framework
- **Testify**: Testing assertions
- **Mockery**: Mock generation
- **go-playground/validator**: Request validation
- **sqlx**: Database operations
- **redis/go-redis**: Redis client
- **simple-migrate**: Database migrations
