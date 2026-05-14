# Project Structure

```
./
├── app
│   ├── apperror      # Global error handler, used in route.go
│   ├── middleware    # Global middleware, used in route.go
│   └── domain/*     # Business logic, separated by domain
├── config            # App config via github.com/caarlos0/env/v11
│                     # Add business config inside config.Business
└── migrations        # SQL migration scripts
```

## Domain Module (`app/domain/*`)

Each domain maps 1-1 to a business concern (e.g. `app/domain/register`). A domain contains:

**`{domain}.go`** — entry point of the domain

- Data models used within the domain (e.g. `User`)
- `New()` function that wires handler and storage together

```go
package example

type User struct {
    FirstName string `json:"firstName"`
    LastName  string `json:"lastName"`
    Age       int    `json:"age"`
}

type domain struct {
    Handler *handler
}

func New() *domain {
    st := NewStorage()
    h := NewHandler(st)
    return &domain{Handler: h}
}
```

**`handler.go`** — declares the handler struct and its dependency interfaces

- Add `//mockery:generate: true` above any interface that needs mocking

```go
package example

//mockery:generate: true
type Storager interface {
    Users() []User
    UserByName(name string) (User, error)
    CreateUser(user User) error
}

type handler struct{
	storage Storager
	timer timer.Timer
}

func NewHandler(storage Storager,timer timer.Timer) *handler {
    return &handler{
    	storage: storage,
     	timer: timer
    }
}
```

**`storage.go`** — database layer using `sqlx`

**`cache.go`** — cache layer (Redis or in-memory)

**`{name}_client.go`** — external API client, named after the 3rd party (e.g. `google_client.go`)

**One file per endpoint**, named after the action:

- `GET /api/v1/user` → `get_user_handler.go`
- `POST /api/v1/user` → `create_user_handler.go`
- `PUT /api/v1/user` → `update_user_handler.go`

## Config

Add new business config fields inside `config.Business`. Use `github.com/caarlos0/env/v11` for env binding.
