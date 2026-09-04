# Middleware

Execution order for HTTP middleware and how context/trace values propagate.

## Global stack (app/echo.go `NewEchoApp`)

Applied in order:

1. `middleware.Recover()`
2. `middleware.CORS("*")`
3. `TagMiddleware()`
4. `RefIDMiddleware(cfg.Header.RefIDKey)`
5. `LoggerMiddleware(cfg.Log.Enable)`

## TagMiddleware

- Reads `ctx.RouteInfo().Name` (the route name passed to `app.GET(echo, "name", ...)`) and sets it as the `tag` in:
  - `ctx` (`ctx.Set(TagKey, name)`)
  - the request context
  - the logger (`ctx.Logger().With(tag, name)`)
- This is why route names must be human-readable.

## RefIDMiddleware

- Reads the ref-id header (default `X-Ref-ID`); generates a UUID if missing.
- Stores it as `traceID` in `ctx`, the request context, and the logger.

## AuthMiddleware

- `app.AuthMiddleware(verifier jwttoken.Verifier)`:
  - Missing `Authorization` header → `Unauthorized(MissingTokenCode, ...)`
  - `Verify` fails due to expiry → `Unauthorized(TokenExpiredCode, ...)`
  - Other verify failure → `Unauthorized(UnauthorizedCode, ...)`
- On success it puts every JWT claim into `ctx` and replaces the request context (`ctx.SetRequest(ctx.Request().WithContext(reqCtx))`), so downstream code sees claims via `c.Get("sub")`.

## ErrorMiddleware

- `app.ErrorMiddleware(handler func(error) app.Error)`:
  - If the handler returns nil → nil.
  - If the error is already an `app.Error` → passes it through unchanged.
  - Otherwise the error goes through the provided mapping (`handleError`), which maps business codes to `app.Error`.

## Logging

- Inside handlers use the echo context logger `ctx.Logger()`, which already carries traceID/tag.
- Every global logging step reads/writes request context so logging stays consistent end-to-end.

## Default error handler

- `errorHandler` in `app/echo.go` converts any non-`app.Error` response into the standard envelope.
- Errors implementing `echo.HTTPStatusCoder` (`*echo.HTTPError`, e.g. 404 route-not-found, 405 method-not-allowed) keep their HTTP status but get business code `HTTP_<status>` (e.g. `"HTTP_404"`).
- Any other unknown error → `InternalErrorCode` ("9999") with HTTP 500.