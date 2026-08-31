# Code Review

Reviewing gotemplate Go code against the conventions in this skill. Go layer by layer over each changed file and flag deviations; when in doubt, read the linked reference. State the rule you are citing so the author can act on it.

## Project structure

- New domain code goes in `internal/{domain}`; HTTP handlers in `app/{domain}app`; reusable infra in `pkg/` (see [architecture.md](architecture.md)).
- No business logic in `main.go` or `app/` common files — they only wire dependencies and register routes.
- `consumer/{domain}consumer` is optional and only for real non-HTTP entrypoints (none in the repo yet).

## Naming & file layout (architecture.md)

- Files connecting to external services are prefixed `adapter_{3party}.go` (`adapter_storage.go`, `adapter_cache.go`, `adapter_client.go`).
- `service.go` holds up to 3–4 exported functions; beyond that, one business per `service_{core}.go` file.
- Helpers shared by many files go in `service_helper.go`, not duplicated.
- Route names are human-readable and use one consistent style (kebab-case, e.g. `get-todos`) — they become the log/metrics `tag`. Flag typos: a route named `helth` would land in every log line.
- Flag misspelled file names — they become imports and are painful to rename later.

## Handlers (app-handler.md)

- Routes registered via `app.GET/POST/PUT/DELETE/PATCH(router, "route_name", path, handler, middlewares...)` inside `RegisterRoute`; `main.go` only calls `NewApp(Deps{...}).RegisterRoute(echo)`.
- `app.ErrorMiddleware(handleError)` on the route group — never per handler.
- Bind + validate in one step with `app.Request[Req](c)`. **Red flag:** plain `c.Bind` next to `validate:"..."` tags — the tags silently do nothing.
- Path params use `param:"id"` and are bound with `c.Bind`.
- Success responses via `app.Ok` (200) / `app.Created` (201); errors are `return`ed and converted by the middleware. Envelope is `{code, success, message, data}` via `app.Response`.
- Auth-protected handlers read the user id with `c.Get("sub").(string)` only behind `app.AuthMiddleware`.
- Handlers stay thin: bind/validate → call service → respond. Business logic belongs in `internal/{domain}`.

## Errors (errors.md)

- 3rd-party/external errors wrapped at the adapter boundary with `serror.From(err)` / `serror.FromWithCode(code, err)`.
- Business codes defined via `serror.NewCoded` inside `internal/{domain}` — no magic strings in service code.
- `handleError` in each `app/{domain}app` maps every business code; catch-all `app.InternalError(app.InternalErrorCode, ...)` covers the rest.
- Business not-found → `app.NotFound` (HTTP 400 by design, not 404).
- Extra context uses `.WithData(...)` (surfaces in response `data`) and `.WithAttr(...)` (log attrs via `logger.ErrorAttrs`).

## Service / storage (architecture.md, errors.md)

- Functions return concrete types, never interfaces; existence checks return `(T, bool, error)` — no pointers for existence.
- Interfaces (`Storager`, `Servicer`) live in the domain package with `//mockery:generate: true` on mockable ones; cross-package dependencies are aliased (`type Hasher = hash.Hasher`).

## Middleware & logging (middleware.md)

- Logging via `ctx.Logger()` (carries traceID/tag) — flag package-level `slog` used inside handlers/services.
- `LoggerMiddleware` reads and restores the request body — flag handlers that read the body again with `io.ReadAll` and don't restore it.
- JWT claims are injected once by `AuthMiddleware`; no per-handler token verification.

## Tests (testing.md)

- Business code has tests for: success path, business-code error (`serror.As` + compare `Code()`/`Msg()`), and dependency error propagation.
- Handler tests: `echotest.ContextConfig{JSONBody: ..., PathValues: ...}.ToContextRecorder(t)`, inject `ctx.Set("sub", ...)`, build the app struct directly (`&todoApp{sv: mock}`).
- Handlers calling `app.Request` need `ctx.Echo().Validator = validator.NewValidator()` in the test, otherwise validation silently passes.
- Storage tests use `go-sqlmock` (`sqlmock.New()` + `sqlx.NewDb(db, "sqlmock")`); mocks live in `mock_test.go` via mockery.
- No `testify/suite` — plain `t.Run` subtests.

## Config / docs (tooling.md)

- New config fields carry `env:"VAR"` + `envDefault` tags and an entry in `.env.example`; remember `{ENV}_VAR` overrides `VAR`.
- New DB schema changes get `migrations/NNNN_*.up.sql`.
- New endpoints annotated with go-swagger (`swagger:route` / `swagger:response` / `swagger:parameters`), success responses wrap `app.SwaggerSuccessResponse`.

## Red flags (quick scan)

- `validate` tags that are never executed (handler uses `c.Bind` instead of `app.Request`).
- Error codes that are neither in `app/const.go` nor a domain `NewCoded`.
- `app.ErrorMiddleware` registered per handler instead of on the group.
- HTTP 404 used for missing business data.
- Functions returning interfaces or `(*T, error)` for existence checks.
- Route names that are unreadable or mix naming styles.
- New env vars missing from `.env.example`.
- New endpoints without swagger annotations or tests.