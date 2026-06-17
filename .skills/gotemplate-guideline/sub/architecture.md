# Architecture

## Domain module structure

Each domain in `app/{domain}/` follows hexagonal (ports-and-adapters) architecture:

```
HTTP → adapter_handler.go (inbound adapter) → Servicer port → service.go (core) → Storager port → adapter_mysql.go (outbound adapter) → DB
```

- `{domain}.go` + `service*.go` — **core**: stdlib-only. Defines models, domain errors, and port interfaces (`Storager`, `Servicer`) with zero third-party imports
- `adapter_handler.go` — **inbound adapter**: translates HTTP (Echo) into calls on the `Servicer` port
- `adapter_mysql.go` — **outbound adapter**: implements `Storager` port via sqlx
- `adapter_module.go` + `main.go` — **composition root**: wires adapters to core at startup (see `app/member/adapter_module.go`)

Third-party imports are confined to `adapter_*.go` files — core files must not import them.

## Wire module

```go
type External struct {
	DB    *sqlx.DB
	Clock clock.Clock
}

type Module struct {
	Handler *handler
}

func NewModule(adp External) *Module {
	st := NewStorage(adp.DB)
	sv := NewService(st, adp.Clock)
	h := NewHandler(sv)
	return &Module{Handler: h}
}
```

Register routes via `Register{Name}Handler(app *app.EchoApp)` on the handler struct.

## Config

Add business env vars inside `config.Config` with `caarlos0/env/v11` tags. Prefix-based override: `ENV=LOCAL` → reads `LOCAL_DATABASE_URL`, etc.
