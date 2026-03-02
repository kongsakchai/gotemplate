package app

import (
	"context"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEchoApp(t *testing.T) {
	t.Run("should start and stop the app without error", func(t *testing.T) {
		w := sync.WaitGroup{}
		w.Add(1)

		app := NewEchoApp()
		go func() {
			defer w.Done()
			err := app.Start(":8080")
			assert.Error(t, http.ErrServerClosed, err)
		}()

		ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancle()
		assert.NoError(t, app.Shutdown(ctx))

		w.Wait()
	})
}

func TestNewMockContext(t *testing.T) {
	t.Run("should create mock context with GET request", func(t *testing.T) {
		method := http.MethodGet
		target := "/api/users"
		payload := ""

		ctx, rec := NewMockContext(method, target, payload)

		assert.NotNil(t, ctx)
		assert.NotNil(t, rec)
		assert.Equal(t, method, ctx.Request().Method)
		assert.Equal(t, target, ctx.Request().URL.Path)
	})

	t.Run("should create mock context with POST request and payload", func(t *testing.T) {
		method := http.MethodPost
		target := "/api/users"
		payload := `{"name":"John","email":"john@example.com"}`

		ctx, rec := NewMockContext(method, target, payload)

		assert.NotNil(t, ctx)
		assert.NotNil(t, rec)
		assert.Equal(t, method, ctx.Request().Method)
		assert.Equal(t, target, ctx.Request().URL.Path)

		body, err := io.ReadAll(ctx.Request().Body)
		assert.NoError(t, err)
		assert.Equal(t, payload, string(body))
	})

	t.Run("should create mock context with PUT request", func(t *testing.T) {
		method := http.MethodPut
		target := "/api/users/1"
		payload := `{"name":"Jane"}`

		ctx, rec := NewMockContext(method, target, payload)

		assert.NotNil(t, ctx)
		assert.NotNil(t, rec)
		assert.Equal(t, method, ctx.Request().Method)
		assert.Equal(t, target, ctx.Request().URL.Path)
	})

	t.Run("should create mock context with DELETE request", func(t *testing.T) {
		method := http.MethodDelete
		target := "/api/users/1"
		payload := ""

		ctx, rec := NewMockContext(method, target, payload)

		assert.NotNil(t, ctx)
		assert.NotNil(t, rec)
		assert.Equal(t, method, ctx.Request().Method)
		assert.Equal(t, target, ctx.Request().URL.Path)
	})

	t.Run("should create response recorder with proper write capabilities", func(t *testing.T) {
		method := http.MethodGet
		target := "/api/users"
		payload := ""

		ctx, rec := NewMockContext(method, target, payload)

		assert.NotNil(t, ctx)
		// Write to response
		rec.WriteHeader(http.StatusOK)
		_, err := rec.Body.WriteString(`{"users":[]}`)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"users":[]}`, rec.Body.String())
	})

	t.Run("should handle URL with query parameters", func(t *testing.T) {
		method := http.MethodGet
		target := "/api/users?page=1&limit=10"
		payload := ""

		ctx, rec := NewMockContext(method, target, payload)

		assert.NotNil(t, ctx)
		assert.NotNil(t, rec)
		assert.Equal(t, "/api/users", ctx.Request().URL.Path)
		assert.Equal(t, "1", ctx.Request().URL.Query().Get("page"))
		assert.Equal(t, "10", ctx.Request().URL.Query().Get("limit"))
	})
}
