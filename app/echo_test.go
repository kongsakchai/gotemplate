package app

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kongsakchai/gotemplate/pkg/config"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type customMarshalerError struct {
	code int
	msg  string
}

func (e *customMarshalerError) Error() string {
	return e.msg
}

func (e *customMarshalerError) StatusCode() int {
	return e.code
}

func (e *customMarshalerError) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.msg)
}

func TestNewEchoApp(t *testing.T) {
	t.Run("should create EchoApp with validator", func(t *testing.T) {
		cfg := config.Config{}
		e := NewEchoApp(cfg)
		assert.NotNil(t, e)
		assert.NotNil(t, e.Validator)
		assert.NotNil(t, e.HTTPErrorHandler)
	})
}

func TestStart(t *testing.T) {
	t.Run("should start and shutdown gracefully", func(t *testing.T) {
		cfg := config.Config{}
		e := NewEchoApp(cfg)
		GET(e, "", "/", func(c *echo.Context) error {
			return c.JSON(200, "ok")
		})
		e.Add(echo.RouteNotFound, "/not-found", func(c *echo.Context) error {
			return c.NoContent(http.StatusNotFound)
		})

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		err := e.Start(ctx, ":0", time.Second)
		assert.NoError(t, err)
	})
}

type failWriter struct {
	http.ResponseWriter
}

func (w *failWriter) Write(p []byte) (n int, err error) {
	return 0, echo.ErrInternalServerError
}

func (w *failWriter) WriteHeader(statusCode int) {
	// do nothing
}

func (w *failWriter) Header() http.Header {
	return http.Header{}
}

func TestErrorHandler(t *testing.T) {
	t.Run("should handle Error", func(t *testing.T) {
		ctx, rec := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/test", nil),
		}.ToContextRecorder(t)

		appErr := BadRequest("4000", "bad request", errors.New("test error"))
		errorHandler(ctx, appErr)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"code":"4000","success":false,"message":"bad request"}`, rec.Body.String())
	})

	t.Run("should handle Error and response fail", func(t *testing.T) {
		ctx, rec := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/test", nil),
		}.ToContextRecorder(t)

		ctx.SetResponse(&failWriter{})

		appErr := BadRequest("4000", "bad request", errors.New("test error"))
		errorHandler(ctx, appErr)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, ``, rec.Body.String())
	})

	t.Run("should handle standard echo error", func(t *testing.T) {
		ctx, rec := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/test", nil),
		}.ToContextRecorder(t)

		echoErr := echo.NewHTTPError(http.StatusNotFound, "not found")
		errorHandler(ctx, echoErr)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.JSONEq(t, `{"code":"HTTP_404","success":false,"message":"not found"}`, rec.Body.String())
	})

	t.Run("should handle unknown error with 500", func(t *testing.T) {
		ctx, rec := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/test", nil),
		}.ToContextRecorder(t)

		errorHandler(ctx, errors.New("unknown error"))

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"code":"9999","success":false,"message":"Internal Server Error"}`, rec.Body.String())
	})

	t.Run("should handle echo HTTPError with empty message", func(t *testing.T) {
		ctx, rec := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/test", nil),
		}.ToContextRecorder(t)

		echoErr := echo.NewHTTPError(http.StatusBadGateway, "")
		errorHandler(ctx, echoErr)

		assert.Equal(t, http.StatusBadGateway, rec.Code)
		assert.JSONEq(t, `{"code":"HTTP_502","success":false,"message":"Bad Gateway"}`, rec.Body.String())
	})

	t.Run("should handle custom marshaler error", func(t *testing.T) {
		ctx, rec := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/test", nil),
		}.ToContextRecorder(t)

		errorHandler(ctx, &customMarshalerError{code: http.StatusTeapot, msg: "custom"})

		assert.Equal(t, http.StatusTeapot, rec.Code)
	})
}

func TestDefaultEchoErrorHandler(t *testing.T) {
	t.Run("should handle nil error gracefully", func(t *testing.T) {
		ctx, rec := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/test", nil),
		}.ToContextRecorder(t)

		defaultEchoErrorHandler(ctx, nil)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("should handle nil error and response fail", func(t *testing.T) {
		ctx, rec := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/test", nil),
		}.ToContextRecorder(t)

		ctx.SetResponse(&failWriter{})

		defaultEchoErrorHandler(ctx, nil)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

// mockRouter implements Router interface for testing
type mockRouter struct {
	LastRoute echo.Route
	AddErr    error
}

func (m *mockRouter) AddRoute(routable echo.Route) (echo.RouteInfo, error) {
	m.LastRoute = routable
	return echo.RouteInfo{}, m.AddErr
}

func TestRouterHelpers(t *testing.T) {
	t.Run("POST registers correct route", func(t *testing.T) {
		mr := &mockRouter{}
		handler := func(c *echo.Context) error { return nil }
		mw := func(next echo.HandlerFunc) echo.HandlerFunc { return next }

		POST(mr, "createUser", "/users", handler, mw)

		assert.Equal(t, http.MethodPost, mr.LastRoute.Method)
		assert.Equal(t, "/users", mr.LastRoute.Path)
		assert.Equal(t, "createUser", mr.LastRoute.Name)
		require.Len(t, mr.LastRoute.Middlewares, 1)
	})

	t.Run("GET registers correct route", func(t *testing.T) {
		mr := &mockRouter{}
		handler := func(c *echo.Context) error { return nil }

		GET(mr, "getUser", "/users/:id", handler)

		assert.Equal(t, http.MethodGet, mr.LastRoute.Method)
		assert.Equal(t, "/users/:id", mr.LastRoute.Path)
		assert.Equal(t, "getUser", mr.LastRoute.Name)
		require.Empty(t, mr.LastRoute.Middlewares)
	})

	t.Run("PUT registers correct route", func(t *testing.T) {
		mr := &mockRouter{}
		handler := func(c *echo.Context) error { return nil }

		PUT(mr, "updateUser", "/users/:id", handler)

		assert.Equal(t, http.MethodPut, mr.LastRoute.Method)
		assert.Equal(t, "/users/:id", mr.LastRoute.Path)
		assert.Equal(t, "updateUser", mr.LastRoute.Name)
	})

	t.Run("DELETE registers correct route", func(t *testing.T) {
		mr := &mockRouter{}
		handler := func(c *echo.Context) error { return nil }

		DELETE(mr, "deleteUser", "/users/:id", handler)

		assert.Equal(t, http.MethodDelete, mr.LastRoute.Method)
		assert.Equal(t, "/users/:id", mr.LastRoute.Path)
		assert.Equal(t, "deleteUser", mr.LastRoute.Name)
	})

	t.Run("PATCH registers correct route", func(t *testing.T) {
		mr := &mockRouter{}
		handler := func(c *echo.Context) error { return nil }

		PATCH(mr, "patchUser", "/users/:id", handler)

		assert.Equal(t, http.MethodPatch, mr.LastRoute.Method)
		assert.Equal(t, "/users/:id", mr.LastRoute.Path)
		assert.Equal(t, "patchUser", mr.LastRoute.Name)
	})

	t.Run("captures AddRoute error from mock", func(t *testing.T) {
		mr := &mockRouter{AddErr: assert.AnError}

		GET(mr, "test", "/test", nil)

		assert.ErrorIs(t, mr.AddErr, assert.AnError)
	})
}
