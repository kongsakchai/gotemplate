# Coding rules

- **Minimize pointers** — use `(T, bool, error)` for existence checks, never `(*T, error)`
- **Wrap errors** with `errs.From(err)` or `errs.New(...)` to preserve stack traces
- **Log through `echo.Context.Logger()`** — never standalone `slog`/`log`
- Compare to zero value for existence: `user == User{}` or add an `IsEmpty()` method
- Use `//mockery:generate: true` above interfaces that need mocking
