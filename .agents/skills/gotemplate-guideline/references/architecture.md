# Architecture

Project layout and naming conventions for `gotemplate`.

## Directory layout

```
internal/{domain}   core/business logic per domain (e.g. internal/auth, internal/todo)
app/{domain}app     HTTP handlers per module (e.g. app/authapp, app/todoapp)
consumer/{domain}consumer   optional — non-HTTP entrypoints per domain; add only when a domain is actually used by both app and consumer (no example in the repo yet)
pkg/                reusable infrastructure (config, logger, serror, validator, jwttoken, hash, httpclient, database, cache, clock, null, migrate)
app/                common code shared across modules (app package)
```

## Naming & file organization inside internal/{domain}

- Name the package after the domain: `internal/todo` package `todo`.
- Files that connect to an external service use the prefix `adapter_{3party}.go`:
  - `adapter_storage.go`, `adapter_cache.go`, `adapter_client.go`
- A service with 3–4 exported functions stays in `service.go`.
- More than that, split one business per file — one business = one function a handler calls, or one endpoint — with prefix `service_{core}.go`:
  - `service_logic.go`, `service_register.go`, `service_get_todo.go`
- Prefix by type so files sort and group by name automatically.
- Split small helper functions inside the file that calls them; keep functions small and easy to test.
- If a test is hard to write, split the function further.
- Functions used by many files go into `service_helper.go`.

## Interface definitions

- Domain interfaces (e.g. `Storager`, `Servicer`) live in the domain package so app and consumer both depend on them.
- Mark interfaces meant for mocking with `//mockery:generate: true`.
- When a struct/interface from another package is needed as dependency, alias it locally:

```go
//mockery:generate: true
type Servicer = auth.Servicer
```

## Shared infra in pkg/

- Anything reused across services or not tied to a specific domain lives in `pkg/`: config, logger, serror, validator, jwttoken, hash, httpclient, database, cache, clock, null, migrate.

## Pointer / boolean rule

- Prefer fewer pointers. If a function needs to signal existence, return a boolean instead (e.g. `(User, bool, error)`).