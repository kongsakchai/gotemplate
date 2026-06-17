# Coding Rules

## Pointer Usage

Minimize pointer usage — avoid entirely when possible.

**When checking if data exists, return a `bool` flag instead of a pointer:**

```go
// ❌ avoid
func (s *storage) UserByName(name string) (*User, error)

// ✅ preferred
func (s *storage) UserByName(name string) (User, bool, error)
```

**Alternatives to pointer for empty checks:**

- If the struct has no slice/map fields, compare directly: `user == User{}`
- Add an `IsEmpty()` method on the struct for reusable empty checks

## Error Handling

**Wrap external library errors** with `errs.From` to preserve stack trace:

```go
import "github.com/kongsakchai/gotemplate/template/common/errs"

row, err := s.db.Query("select * from user")
if err == sql.ErrNoRows {
    return User{}, false, nil
} else if err != nil {
    return User{}, false, errs.From(err)
}
```

**Create custom errors** with `errs.New` to preserve stack trace:

```go
return User{}, false, errs.New("user %s not found", name)
```

## Logging

Use the logger from `echo.Context` to ensure logs include tag and trace ID:

```go
ctx.Logger().InfoContext(ctx, msg, args...)
```

> Do not use a standalone logger (e.g. `log.Println`, `slog`) directly — always log through `echo.Context`.
