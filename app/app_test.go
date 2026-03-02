package app

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAppResponse(t *testing.T) {
	t.Run("should return 200 OK when use Ok", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		expectedStatus := http.StatusOK
		expectedResp := "{\"code\":\"0000\",\"success\":true,\"data\":\"Success\",\"message\":\"Ping!\"}\n"

		Ok(ctx, "Success", "Ping!")

		assert.Equal(t, expectedStatus, rec.Code)
		assert.JSONEq(t, expectedResp, rec.Body.String())
	})

	t.Run("should return 201 Created when use Created", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		expectedStatus := http.StatusCreated
		expectedResp := "{\"code\":\"0000\",\"success\":true,\"data\":\"Created\",\"message\":\"Ping!\"}\n"

		Created(ctx, "Created", "Ping!")

		assert.Equal(t, expectedStatus, rec.Code)
		assert.JSONEq(t, expectedResp, rec.Body.String())
	})

	t.Run("should return 500 Internal Server Error when use FailWithError", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		expectedStatus := http.StatusInternalServerError
		expectedResp := "{\"code\":\"5000\",\"success\":false,\"message\":\"Internal Server Error\"}\n"

		err := InternalError("5000", "Internal Server Error", errors.New("unexpected error"))
		Fail(ctx, err)

		assert.Equal(t, expectedStatus, rec.Code)
		assert.JSONEq(t, expectedResp, rec.Body.String())
	})
}
