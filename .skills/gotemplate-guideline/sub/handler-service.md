# Handler & service conventions

## Handler

- Handler struct holds a `Servicer` (not storage directly)
- Handle (exported) → validates input via `app.Request(ctx, &req)`, calls service, returns response
- Error mapping via `handlerError()` — switches on `errors.Is(err, DomainError)`:
  - Matched domain error → return specific `app.BadRequest` / `app.Conflict` / `app.NotFound` (with business code)
  - Default → return `app.InternalError` for unexpected errors
- Use `app.Request(ctx, &req)` for bind + validate in one call
- **`app.NotFound()` returns HTTP 400** — 404 is reserved for missing routes only

## Service

- `service.go` — `service` struct + `NewService(storage Storager, clock Clock)`
- `{domain}.go` — `Servicer` interface defining business logic contract
- One file per service method: `service_create.go`, `service_get_all.go`, `service_get_by_username.go`
- Service adds business logic (validation, domain errors) on top of storage calls
- Storage returns `(T, bool, error)`; service converts to `(T, error)` with domain errors like `ErrorMemberNotFound`

## API responses

| Helper | Status |
|--------|--------|
| `app.Ok(ctx, data, msg...)` | 200 |
| `app.Created(ctx, data, msg...)` | 201 |
| `app.Fail(ctx, error)` | varies |
| `app.BadRequest(code, msg, err)` | 400 |
| `app.NotFound(code, msg, err)` | 400 ⚠️ |
| `app.Unauthorized(code, msg, err)` | 401 |
| `app.Forbidden(code, msg, err)` | 403 |
| `app.Conflict(code, msg, err)` | 409 |
| `app.InternalError(code, msg, err)` | 500 |
