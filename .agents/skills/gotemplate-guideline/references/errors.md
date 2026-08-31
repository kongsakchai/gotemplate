# Errors

How errors flow through the project: `serror` in business code, `app.Error` + `handleError` at the HTTP boundary.

## serror (pkg/serror)

Tracing wrapper for errors. Produces a message like:

```
error: <msg>, code: <code>, at: (file.go:line) package.func
```

Use cases:

- Define a business error code: `serror.NewCoded("code", "message")`.
  - Call `.Err()` to produce the error: `return ErrTodoNotFound.Err()`.
- Wrap an external/3rd-party error to keep it traceable: `serror.From(err)`.
- Wrap with a code: `serror.NewWithCode(code, "...")`, `serror.FromWithCode(code, err)`.
- Attach extra context, chainable:
  - `.WithData(...)` → surfaced as `app.Error.Data` in the HTTP response.
  - `.WithAttr(...)` → included in log output via `logger.ErrorAttrs(err)`.
- Check if an error is a `serror`: `e, ok := serror.As(err)`.
- Errors compare with `errors.Is`/`errors.As`; `e.Code()`, `e.Msg()` give code and message.

Wrap external errors (DB, HTTP client, driver) with `serror.From` at the adapter boundary so the origin is tracked.

## app.Error (app/error.go)

HTTP-facing error type:

```go
type Error struct {
	HTTPCode int
	Code     string
	Message  string
	Data     any
	Err      error
}
```

Constructors (code, msg, underlying err, data...):

- `app.InternalError` → 500
- `app.BadRequest` → 400
- `app.NotFound` → **400 on purpose** (see rule below)
- `app.Unauthorized` → 401
- `app.Forbidden` → 403
- `app.Conflict` → 409

### NotFound uses 400, not 404

400 signals "the request input is wrong so the data can't be found"; it avoids callers having to disambiguate "data missing" from "route missing". Keep this convention.

## Global business codes (app/const.go)

| Code | Meaning |
|------|---------|
| `0000` | success |
| `1000` | bad request |
| `1001` | invalid request (validation) |
| `10002` | missing token |
| `10003` | unauthorized |
| `10004` | token expired |
| `10005` | invalid token |
| `9998` | database not ready |
| `9999` | internal error |

Module-level business codes are defined in each `internal/{domain}` package via `serror.NewCoded`.

## Per-module handleError

Every `app/{domain}app` maps business codes → `app.Error` with its own `handleError`:

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

The catch-all (`app.InternalError`) covers unknown errors and typed-nil errors.

## Wiring

- `app.ErrorMiddleware(handleError)` is registered once on the route group in `RegisterRoute` — never per handler.
- Handlers just `return err`; middleware converts it.
- The default error handler in `app/echo.go` wraps any non-`app.Error` response in the standard envelope: `echo.HTTPStatusCoder` errors (e.g. 404 route not found, 405 method not allowed) keep their HTTP status but get business code `HTTP_<status>`; anything else falls back to `InternalErrorCode` ("9999") / HTTP 500.