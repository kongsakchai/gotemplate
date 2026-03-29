package apperror

import (
	"errors"
	"net/http"
	"testing"

	"github.com/kongsakchai/gotemplate/template/app"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
)

type errorWriter struct {
	http.ResponseWriter
}

func (w *errorWriter) Write(p []byte) (n int, err error) {
	return 0, echo.ErrInternalServerError
}

func (w *errorWriter) WriteHeader(statusCode int) {
	// do nothing
}

func (w *errorWriter) Header() http.Header {
	return http.Header{}
}

type jsonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e jsonError) MarshalJSON() ([]byte, error) {
	return []byte(`{"message":"error"}`), nil
}

func (e jsonError) Error() string {
	return e.Message
}

func (e jsonError) StatusCode() int {
	return e.Code
}

func TestErrorHandler(t *testing.T) {
	t.Run("should return bad request when error is app error", func(t *testing.T) {
		ctx, rec := echotest.ContextConfig{}.ToContextRecorder(t)

		appErr := app.BadRequest("4000", "error", errors.New("error"))
		ErrorHandler(ctx, appErr)

		assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
	})

	t.Run("should return internal error when error isn't app error", func(t *testing.T) {
		ctx, rec := echotest.ContextConfig{}.ToContextRecorder(t)

		ErrorHandler(ctx, errors.New("error"))

		assert.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
	})

	t.Run("should return http error when error is http error", func(t *testing.T) {
		ctx, rec := echotest.ContextConfig{}.ToContextRecorder(t)

		ErrorHandler(ctx, &echo.HTTPError{Code: http.StatusNotFound, Message: "bad request"})

		assert.Equal(t, http.StatusNotFound, rec.Result().StatusCode)
	})

	t.Run("should return http error when error is http error without message", func(t *testing.T) {
		ctx, rec := echotest.ContextConfig{}.ToContextRecorder(t)

		ErrorHandler(ctx, &echo.HTTPError{Code: http.StatusNotFound, Message: ""})

		assert.Equal(t, http.StatusNotFound, rec.Result().StatusCode)
	})

	t.Run("should return error when error is json Marshaler", func(t *testing.T) {
		ctx, rec := echotest.ContextConfig{}.ToContextRecorder(t)

		ErrorHandler(ctx, jsonError{Code: http.StatusInternalServerError, Message: "error"})

		assert.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
	})

	t.Run("should no response when send response fail", func(t *testing.T) {
		ctx := echotest.ContextConfig{}.ToContext(t)
		w := &errorWriter{}
		ctx.SetResponse(w)

		appErr := app.BadRequest("4000", "error", errors.New("error"))
		ErrorHandler(ctx, appErr)
	})
}
