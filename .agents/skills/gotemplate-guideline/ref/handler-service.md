# Handler & service conventions

## Handler

- Handler struct holds a `Servicer` (not storage directly)
- Handle (exported) → validates input via `app.Request(ctx, &req)`, calls service, returns response
- Error mapping via `handlerError()` — switches on `errors.Is(err, DomainError)`:
    - Matched domain error → return specific `app.BadRequest` / `app.Conflict` / `app.NotFound` (with business code)
    - Default → return `app.InternalError` for unexpected errors
- Use `app.Request(ctx, &req)` for bind + validate in one call
- **`app.NotFound()` returns HTTP 400** — 404 is reserved for missing routes only
- wires adapters to core at startup.

## Wire module

```go
// app/memberapp/handler_member.go
type Deps struct {
	DB    *sqlx.DB
	Clock someDomain.Clock
}

func NewMemberHandler(deps Deps) *memberHandler {
		st := member.NewStorage(deps.DB)
		sv := member.NewService(member.ServiceDeps{
			Clock:   deps.Clock,
			Storage: st,
		})

		return &memberHandler{
				service:    sv,
		}
}

func (h *memberHandler) RegisterHandler(e *echo.Echo) {
	apiPath := e.Group("/api/v1/members")
	apiPath.GET("", h.members, app.TagMiddleware("get_members"))
}

func (h *memberHandler) handlerError(err error) error {
	if ... {
		return app.BadRequest(d.Code, d.Message, err)
	}
	return app.InternalError(app.InternalErrorCode, app.InternalErrorMsg, err)
}

func (h *memberHandler) members(ctx *echo.Context) error {
	members, err := h.service.GetMembers(ctx.Request().Context())
	if err != nil {
		return h.handlerError(err)
	}
	return app.Ok(ctx, members)
}
```

## Service

- `service.go` — `service` struct + `NewService(storage Storager, clock Clock)`
- `{domain}.go` — `Servicer` interface defining business logic contract
- One file per service method: `service_create.go`, `service_get_all.go`, `service_get_by_username.go`
- Service adds business logic (validation, domain errors) on top of storage calls
- Storage returns `(T, bool, error)`; service converts to `(T, error)` with domain errors like `ErrorMemberNotFound`

## API responses

| Helper                              | Status |
| ----------------------------------- | ------ |
| `app.Ok(ctx, data, msg...)`         | 200    |
| `app.Created(ctx, data, msg...)`    | 201    |
| `app.Fail(ctx, error)`              | varies |
| `app.BadRequest(code, msg, err)`    | 400    |
| `app.NotFound(code, msg, err)`      | 400 ⚠️ |
| `app.Unauthorized(code, msg, err)`  | 401    |
| `app.Forbidden(code, msg, err)`     | 403    |
| `app.Conflict(code, msg, err)`      | 409    |
| `app.InternalError(code, msg, err)` | 500    |
