# App / Handler

Building the HTTP layer in `app/{domain}app`.

## Module setup

- Handlers live in `app/{domain}app` so each module is isolated and can grow independently.
- If a domain later needs a non-HTTP entrypoint (worker, job consumer), optionally add `consumer/{domain}consumer` — both depend on `internal/{domain}` without carrying each other's code. This is a forward-looking convention; the repo has no consumer yet.
- `app/{domain}app/{domain}app.go` exposes `NewApp(Deps{...})` and `RegisterRoute(echo *app.EchoApp)` so `main` can be split into microservices later.

```go
type Deps struct {
	DB       *sqlx.DB
	Verifier jwttoken.Verifier
}

func NewApp(deps Deps) *todoApp {
	st := todo.NewStorage(deps.DB)
	sv := todo.NewService(st)
	return &todoApp{sv: sv, verifier: deps.Verifier}
}

func (a *todoApp) RegisterRoute(echo *app.EchoApp) {
	g := echo.Group("/api/v1/todos", app.AuthMiddleware(a.verifier), app.ErrorMiddleware(handleError))
	app.GET(g, "get-todos", "", a.getTodos)
	app.POST(g, "create-todo", "", a.createTodo)
	app.PUT(g, "update-todo", "/:id", a.updateTodo)
	app.DELETE(g, "delete-todo", "/:id", a.deleteTodo)
}
```

- Keep one consistent name (`RegisterRoute` in every domain app; `todoapp` legacy uses `RegisterRoutes` — prefer `RegisterRoute`).

## Routing

- Register routes with the helpers in `app/echo_route.go`: `app.GET`, `app.POST`, `app.PUT`, `app.DELETE`, `app.PATCH`.
- Signature is `app.METHOD(router, "route_name", path, handler, middlewares...)`.
- Route names must be readable — they become the `tag` used in logs/metrics.

## Request / validation

- Bind + validate in one step with `app.Request[Req](c)`; it returns a bad-request error using `InValidCode` ("1001") on `Validate` failure and `BadRequestCode` ("1000") on bind failure.
- Define request structs with `json` + `validate` tags (e.g. `validate:"required"`).
- Path params use `param:"id"` and are bound with `c.Bind(&req)`.
- `app.Request` returns `app.Error`, so let it bubble up through `ErrorMiddleware`.

```go
type createTodoRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
```

## Responses

- Success responses use `app.Ok(ctx, data)` (200) and `app.Created(ctx, data)` (201); errors use `app.Fail(ctx, appErr)`.
- Response envelope (see `app/app.go`):

```json
{ "code": "0000", "success": true, "message": "", "data": {} }
```

## Auth in handlers

- `app.AuthMiddleware(verifier)` verifies JWT and puts every claim into the context; read the user id with `c.Get("sub").(string)`.