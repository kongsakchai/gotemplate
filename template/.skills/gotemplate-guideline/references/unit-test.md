# Unit Test Convention

## File & Package

- Name the test file after the source file: `get_user_handler.go` → `get_user_handler_test.go`
- Use the **same package** as the source file to access unexported functions

## Mocking with Mockery

Add `//mockery:generate: true` above any interface that needs a mock:

```go
//mockery:generate: true
type Storager interface {
    Users() []User
    UserByName(name string) (User, error)
    CreateUser(user User) error
}
```

Update `.mockery` to point to the correct project path:

```yaml
# .mockery
packages:
    github.com/your-org/your-project/app/domain/example: ...
```

Then generate mocks:

```sh
mockery
# or
make genmock
```

## Common Mock Utilities

For shared/common mocks, use the `mockutil` package:

```go
import "github.com/kongsakchai/gotemplate/template/common/pkg/mockutil"

mockutil.NewTimer(t)  // mock for timer
```

Use `mockutil` for any common dependency that is shared across multiple domains rather than generating a new mock per domain.

## Writing Tests

```go
package example

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/kongsakchai/gotemplate/template/common/pkg/mockutil"
    "github.com/your-org/your-project/app/domain/example/mocks"
)

func TestGetUser(t *testing.T) {
    tests := []struct {
        name      string
        mockSetup func(*mocks.Storager)
        input     string
        wantErr   bool
    }{
        {
            name: "success",
            mockSetup: func(m *mocks.Storager) {
                m.On("UserByName", "john").Return(User{FirstName: "john"}, nil)
            },
            input:   "john",
            wantErr: false,
        },
        {
            name: "user not found",
            mockSetup: func(m *mocks.Storager) {
                m.On("UserByName", "unknown").Return(User{}, nil)
            },
            input:   "unknown",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := mocks.NewStorager(t)
            tt.mockSetup(m)
            h := NewHandler(m)

            result, err := h.getUser(tt.input)
            if tt.wantErr {
                assert.False(t, err.IsEmpty())
            } else {
                assert.True(t, err.IsEmpty())
                assert.Equal(t, tt.input, result.FirstName)
            }
        })
    }
}
```
