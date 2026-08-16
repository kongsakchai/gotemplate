# Architecture

## Domain module structure

Each domain in `internal/{domain}/` follows hexagonal (ports-and-adapters) architecture:

```
HTTP → handler.go (inbound adapter) → Servicer port → service.go (core) → Storager port → adapter_mysql.go (outbound adapter) → DB
```

- `{domain}.go` + `service*.go` — **core**: Defines models, domain errors, and port interfaces (`Storager`, `Servicer`)
- `adapter_mysql.go` — **outbound adapter**: implements `Storager` port via sqlx
- `adapter_external.go` - **outbound adapter**: implements `Externa;` port via httpclient
- `adapter_cache.go` - **outbound adapter**: implements `Cache` port via redis, local cache
- `app/{domain}app/handler_{domain}.go` - **inbound adapter**: Translates HTTP (Echo) into calls on the `Servicer` port and map error from service to http response.

Third-party imports are confined to `adapter_*.go` files — core files must not import them.

Register routes via `RegisterHandler(app *app.EchoApp)` on the handler struct.

## Config

Add business env vars inside `config.Config` with `caarlos0/env/v11` tags. Prefix-based override: `ENV=LOCAL` → reads `LOCAL_DATABASE_URL`, etc.
