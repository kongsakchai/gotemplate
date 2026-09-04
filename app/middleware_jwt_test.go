package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
)

type fakeVerifier struct {
	verify func(token string) (map[string]any, error)
}

func (f fakeVerifier) Verify(token string) (map[string]any, error) {
	return f.verify(token)
}

func TestAuthMiddleware(t *testing.T) {
	newCtx := func(authHeader string) (*echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		if authHeader != "" {
			req.Header.Set(authorizationHeader, authHeader)
		}
		ctx, rec := echotest.ContextConfig{
			Request: req,
		}.ToContextRecorder(t)
		return ctx, rec
	}

	t.Run("missing token returns unauthorized", func(t *testing.T) {
		ctx, _ := newCtx("")

		verifierCalled := false
		nextCalled := false
		handler := AuthMiddleware(fakeVerifier{
			verify: func(token string) (map[string]any, error) {
				verifierCalled = true
				return nil, nil
			},
		})(func(c *echo.Context) error {
			nextCalled = true
			return nil
		})

		err := handler(ctx)

		assert.False(t, verifierCalled)
		assert.False(t, nextCalled)

		var appErr Error
		assert.ErrorAs(t, err, &appErr)
		assert.Equal(t, http.StatusUnauthorized, appErr.HTTPCode)
		assert.Equal(t, MissingTokenCode, appErr.Code)
		assert.Equal(t, MissingTokenMsg, appErr.Message)
	})

	t.Run("token shorter than Bearer prefix returns invalid token", func(t *testing.T) {
		ctx, _ := newCtx("Bearer")

		nextCalled := false
		handler := AuthMiddleware(fakeVerifier{
			verify: func(token string) (map[string]any, error) {
				t.Fatal("verifier should not be called")
				return nil, nil
			},
		})(func(c *echo.Context) error {
			nextCalled = true
			return nil
		})

		err := handler(ctx)

		assert.False(t, nextCalled)

		var appErr Error
		assert.ErrorAs(t, err, &appErr)
		assert.Equal(t, http.StatusUnauthorized, appErr.HTTPCode)
		assert.Equal(t, InvalidTokenCode, appErr.Code)
		assert.Equal(t, InvalidTokenMsg, appErr.Message)
	})

	t.Run("expired token returns token expired", func(t *testing.T) {
		ctx, _ := newCtx("Bearer expired-token")

		nextCalled := false
		handler := AuthMiddleware(fakeVerifier{
			verify: func(token string) (map[string]any, error) {
				return nil, jwt.ErrTokenExpired
			},
		})(func(c *echo.Context) error {
			nextCalled = true
			return nil
		})

		err := handler(ctx)

		assert.False(t, nextCalled)

		var appErr Error
		assert.ErrorAs(t, err, &appErr)
		assert.Equal(t, http.StatusUnauthorized, appErr.HTTPCode)
		assert.Equal(t, TokenExpiredCode, appErr.Code)
		assert.Equal(t, TokenExpiredMsg, appErr.Message)
		assert.ErrorIs(t, appErr.Err, jwt.ErrTokenExpired)
	})

	t.Run("verify error returns unauthorized", func(t *testing.T) {
		ctx, _ := newCtx("Bearer bad-token")

		nextCalled := false
		handler := AuthMiddleware(fakeVerifier{
			verify: func(token string) (map[string]any, error) {
				return nil, assert.AnError
			},
		})(func(c *echo.Context) error {
			nextCalled = true
			return nil
		})

		err := handler(ctx)

		assert.False(t, nextCalled)

		var appErr Error
		assert.ErrorAs(t, err, &appErr)
		assert.Equal(t, http.StatusUnauthorized, appErr.HTTPCode)
		assert.Equal(t, UnauthorizedCode, appErr.Code)
		assert.Equal(t, UnauthorizedMsg, appErr.Message)
	})

	t.Run("valid token sets claims and calls next", func(t *testing.T) {
		ctx, _ := newCtx("Bearer valid-token")

		var gotToken string
		nextCalled := false
		handler := AuthMiddleware(fakeVerifier{
			verify: func(token string) (map[string]any, error) {
				gotToken = token
				return map[string]any{"sub": "user-1", "email": "test@example.com"}, nil
			},
		})(func(c *echo.Context) error {
			nextCalled = true
			assert.Equal(t, "user-1", c.Get("sub"))
			assert.Equal(t, "test@example.com", c.Get("email"))
			assert.Equal(t, "user-1", c.Request().Context().Value("sub"))
			return nil
		})

		err := handler(ctx)

		assert.NoError(t, err)
		assert.True(t, nextCalled)
		assert.Equal(t, "valid-token", gotToken)
		assert.Equal(t, "user-1", ctx.Get("sub"))
		assert.Equal(t, "test@example.com", ctx.Get("email"))
		assert.Equal(t, "user-1", ctx.Request().Context().Value("sub"))
	})
}
