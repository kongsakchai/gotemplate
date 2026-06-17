package app

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kongsakchai/gotemplate/template/pkg/config"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
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
		e.GET("/", func(c *echo.Context) error {
			return c.JSON(200, "ok")
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
	t.Run("should handle app.Error", func(t *testing.T) {
		ctx, rec := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/test", nil),
		}.ToContextRecorder(t)

		appErr := BadRequest("4000", "bad request", errors.New("test error"))
		errorHandler(ctx, appErr)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"code":"4000","success":false,"message":"bad request"}`, rec.Body.String())
	})

	t.Run("should handle app.Error and response fail", func(t *testing.T) {
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
		assert.JSONEq(t, `{"code":"","success":false,"message":"not found"}`, rec.Body.String())
	})

	t.Run("should handle unknown error with 500", func(t *testing.T) {
		ctx, rec := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/test", nil),
		}.ToContextRecorder(t)

		errorHandler(ctx, errors.New("unknown error"))

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"code":"","success":false,"message":"Internal Server Error"}`, rec.Body.String())
	})

	t.Run("should handle echo HTTPError with empty message", func(t *testing.T) {
		ctx, rec := echotest.ContextConfig{
			Request: httptest.NewRequest(http.MethodGet, "/test", nil),
		}.ToContextRecorder(t)

		echoErr := echo.NewHTTPError(http.StatusBadGateway, "")
		errorHandler(ctx, echoErr)

		assert.Equal(t, http.StatusBadGateway, rec.Code)
		assert.JSONEq(t, `{"code":"","success":false,"message":"Bad Gateway"}`, rec.Body.String())
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
