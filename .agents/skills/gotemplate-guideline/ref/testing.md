# Test conventions

- Use **same package** as source to access unexported functions
- **Mock Echo context** with `github.com/labstack/echo/v5/echotest`
- **Mock interfaces** with mockery (generated in-source as `mock_*_test.go`)
- **DB tests**: in-memory SQLite (`modernc.org/sqlite`) or `github.com/DATA-DOG/go-sqlmock`
- **Assertions**: `github.com/stretchr/testify` (assert/require)
