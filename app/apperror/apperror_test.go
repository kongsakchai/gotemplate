package apperror

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockFailContext struct {
	echo.Context
}

func (m mockFailContext) JSON(code int, i any) error {
	return errors.New("JSON ERROR")
}

func (m mockFailContext) Request() *http.Request {
	return &http.Request{}
}

func (m mockFailContext) Get(key string) any {
	return key
}

func TestErrorHandler(t *testing.T) {
	t.Run("should return bad request when error is app error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		appErr := app.BadRequest("4000", "error", errors.New("error"))
		ErrorHandler(appErr, ctx)

		assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
	})

	t.Run("should return internal error when error isn't app error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		ErrorHandler(errors.New("error"), ctx)

		assert.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
	})

	t.Run("should return http error when error is http error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		ErrorHandler(&echo.HTTPError{Code: http.StatusNotFound, Message: "bad request"}, ctx)

		assert.Equal(t, http.StatusNotFound, rec.Result().StatusCode)
	})

	t.Run("should no response when send response fail", func(t *testing.T) {
		ctx := &mockFailContext{}

		appErr := app.BadRequest("4000", "error", errors.New("error"))
		ErrorHandler(appErr, ctx)
	})
}
