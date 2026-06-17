# Handler Convention

## Structure

- **One file per endpoint**, named after the action: `get_user_handler.go`, `update_user_handler.go`
- **No separate Service layer** — handler and service logic live in the same file
    - `Handle` function (exported, PascalCase) — receives `*echo.Context`, validates input, calls process, returns response
    - `Process` function (unexported, camelCase) — contains business logic, returns data + `app.Error`
    - Process functions can be broken down into smaller unexported sub-functions to:
        - Make each part independently unit-testable
        - Share logic across multiple endpoints within the same domain
- Request struct (optional)
- Response struct (optional)

```go
type CreateUserRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Age       int    `json:"age" validate:"required,gte=0,lte=130"`
}

func (h *handler) CreateUser(ctx echo.Context) error {
	var req CreateUserRequest
	if err := ctx.Bind(&req); err != nil {
		return app.BadRequest("4001", "invalid request body", err)
	}
	if err := ctx.Validate(&req); err != nil {
		return app.BadRequest("4002", "validation error", err, err)
	}
	if err := h.createUser(req); !err.IsEmpty() {
		return err
	}
	return app.Create(ctx, nil)
}

func (h *handler) createUser(req CreateUserRequest) app.Error {
	if err := h.validateUserNotExists(req.FirstName); !err.IsEmpty() {
		return err
	}
	if err := h.storage.CreateUser(User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Age:       req.Age,
	}); err != nil {
		return app.InternalServer("5002", "failed to create user", err)
	}
	return app.Error{}
}

func (h *handler) validateUserNotExists(firstName string) app.Error {
	user, err := h.storage.UserByName(firstName)
	if err != nil {
		return app.InternalServer("5001", "failed to get user by name", err)
	}
	if user.FirstName != "" {
		return app.Conflict("4003", "user already exists", nil)
	}
	return app.Error{}
}
```

## API Response

**Success**
| Function | Status |
|---|---|
| `app.Ok(ctx, data)` | 200 |
| `app.Create(ctx, data)` | 201 |

**Error** — return `app.Error` or `error`, handled by global error handler in `route.go`

| Function                                      | Status |
| --------------------------------------------- | ------ |
| `app.BadRequest(code, msg, err, data...)`     | 400    |
| `app.Unauthorized(code, msg, err, data...)`   | 401    |
| `app.Forbidden(code, msg, err, data...)`      | 403    |
| `app.NotFound(code, msg, err, data...)`       | 400 ⚠️ |
| `app.Conflict(code, msg, err, data...)`       | 409    |
| `app.InternalServer(code, msg, err, data...)` | 500    |

> ⚠️ **Data not found uses `400` not `404`** — `404` is reserved for missing routes/endpoints.
> Missing data is treated as an invalid client request, not a missing resource.
