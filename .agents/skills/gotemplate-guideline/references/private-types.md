# Private Concrete Types

In `gotemplate`, every concrete implementation is **unexported (private)**. The package exports only the *contract*: domain interfaces, deps structs, and domain models. Consumers never name the concrete type — they receive it from a constructor and use it through its exported methods.

## The pattern at a glance

| What | Naming | Visibility | Example location |
|------|--------|------------|------------------|
| Service implementation | `service` | private | `internal/auth/service.go` |
| Storage adapter | `storage` | private | `internal/auth/adapter_storage.go` |
| App / handler implementation | `{domain}app` | private | `app/authapp/authapp.go` |
| DB record / row struct | `{name}Record` | private | `userRecord`, `todoRecord` |
| Constructor args | `ServiceDeps` / `Deps` | exported | `auth.ServiceDeps`, `authapp.Deps` |
| Contract interfaces | `Servicer`, `Storager` | exported | `internal/auth/auth.go` |
| Domain models | `User`, `Todo` | exported | `auth.User`, `todo.Todo` |

## Rules

1. **Declare the concrete type unexported** — lowercase name:
   ```go
   type service struct { st Storager }
   type storage struct { db *sqlx.DB }
   type authApp struct { sv Servicer }
   ```
2. **Constructors return the private concrete type, never the interface**:
   ```go
   func NewService(deps ServiceDeps) *service
   func NewStorage(db *sqlx.DB) *storage
   func NewApp(deps Deps) *authApp
   ```
   This matches the "return concrete types, never interfaces" rule.
3. **Export the contract as interfaces in the domain file** (`auth.go` / `todo.go`): `Servicer`, `Storager`. Mark mock-ready interfaces with `//mockery:generate: true`.
4. **Constructor args go in an exported deps struct**: `ServiceDeps` in `internal/{domain}`, `Deps` in `app/{domain}app`.
5. **Fields inside the private struct are typed as the interface**, not the concrete type:
   ```go
   type service struct { st Storager }
   type authApp struct { sv Servicer }
   ```
6. **Alias external dependency types locally** so the package owns its contract:
   ```go
   //mockery:generate: true
   type Hasher = hash.Hasher
   type Signer = jwttoken.Signer
   type Verifier = jwttoken.Verifier
   ```
7. **Private structs cross package boundaries through constructors** — it is legal Go to call `authapp.NewApp(...)` in `main.go` and use the returned value's exported methods only (e.g. `RegisterRoute(echo)`). You never need the concrete type name outside its package.

## Example (inline copy from auth module)

Copied inline so this skill stays valid even if the example modules are removed from the repo. The `todo` module follows the same shape.

### 1. Domain contract — `internal/{domain}/{domain}.go`

```go
package auth

import (
	"context"
	"time"
	"github.com/kongsakchai/gotemplate/pkg/hash"
	"github.com/kongsakchai/gotemplate/pkg/jwttoken"
)

var ErrUserAlreadyExists = serror.NewCoded("2000", "user already exists")

type User struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

//mockery:generate: true
type Storager interface {
	CreateUser(ctx context.Context, user User) error
	FindUserByUsername(ctx context.Context, username string) (User, bool, error)
}

//mockery:generate: true
type Hasher = hash.Hasher

//mockery:generate: true
type Signer = jwttoken.Signer

type Servicer interface {
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, username, password string) error
}
```

### 2. Service impl — `internal/{domain}/service.go`

```go
package auth

type service struct {
	st     Storager
	hasher Hasher
	signer Signer
}

type ServiceDeps struct {
	Storager Storager
	Hasher   hash.Hasher
	Signer   jwttoken.Signer
}

func NewService(deps ServiceDeps) *service {
	return &service{
		st:     deps.Storager,
		hasher: deps.Hasher,
		signer: deps.Signer,
	}
}

func (s *service) Register(ctx context.Context, username, password string) error {
	_, found, err := s.st.FindUserByUsername(ctx, username)
	if err != nil {
		return err
	}
	if found {
		return ErrUserAlreadyExists.Err()
	}
	hashedPassword, _ := s.hasher.HashPassword(password)
	return s.st.CreateUser(ctx, User{
		Id:       uuid.NewString(),
		Username: username,
		Password: hashedPassword,
	})
}
```

### 3. Storage adapter — `internal/{domain}/adapter_storage.go`

```go
package auth

type storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *storage {
	return &storage{db: db}
}

// Private row struct mapped to the exported domain model.
type userRecord struct {
	Id        string    `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password_hash"`
	CreatedAt time.Time `db:"created_at"`
}

func (r userRecord) toUser() User {
	return User{Id: r.Id, Username: r.Username, Password: r.Password, CreatedAt: r.CreatedAt}
}

func (s *storage) FindUserByUsername(ctx context.Context, username string) (User, bool, error) {
	var rec userRecord
	err := s.db.GetContext(ctx, &rec, "SELECT * FROM user WHERE username = ?", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, false, nil
		}
		return User{}, false, serror.From(err)
	}
	return rec.toUser(), true, nil
}
```

### 4. App handler — `app/{domain}app/{domain}app.go`

```go
package authapp

import (
	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/internal/auth"
	"github.com/kongsakchai/gotemplate/pkg/hash"
	"github.com/kongsakchai/gotemplate/pkg/jwttoken"
)

type authApp struct {
	sv Servicer
}

// Type alias for external dependencies so the app package owns its contract.
type Servicer = auth.Servicer
type Hasher = hash.Hasher
type Signer = jwttoken.Signer

type Deps struct {
	DB     *sqlx.DB
	Hasher Hasher
	Signer Signer
}

func NewApp(deps Deps) *authApp {
	st := auth.NewStorage(deps.DB)
	sv := auth.NewService(auth.ServiceDeps{
		Storager: st,
		Hasher:   deps.Hasher,
		Signer:   deps.Signer,
	})
	return &authApp{sv: sv}
}

func (a *authApp) RegisterRoute(echo *app.EchoApp) {
	g := echo.Group("/api/v1/auth", app.ErrorMiddleware(handleError))
	app.POST(g, "auth_register", "/register", a.register)
	app.POST(g, "auth_login", "/login", a.login)
}
```

### Usage outside the package

```go
// main.go — unexported type is used only through its constructor return value.
authapp.NewApp(authapp.Deps{DB: db, Hasher: hasher, Signer: jwt}).RegisterRoute(echo)
```

## Why

- The concrete type is an implementation detail; exporting it would turn every internal change into a public API change.
- A single exported contract (interface) per collaborator keeps the app layer, tests, and mocks coupled to behavior, not to internals.
- Combined with mockery (`//mockery:generate: true`), this makes unit tests trivial: app tests mock `Servicer`, service tests mock `Storager`.
