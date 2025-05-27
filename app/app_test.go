package app

import (
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
		expectedResp := "{\"code\":\"0000\",\"status\":\"SUCCESS\",\"message\":\"\",\"data\":\"Success\"}\n"

		Ok(ctx, "Success")

		assert.Equal(t, expectedStatus, rec.Code)
		assert.Equal(t, expectedResp, rec.Body.String())
	})

	t.Run("should return 200 OK with message when use OkWithMessage", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		expectedStatus := http.StatusOK
		expectedResp := "{\"code\":\"0000\",\"status\":\"SUCCESS\",\"message\":\"Pong!\",\"data\":\"Success\"}\n"

		OkWithMessage(ctx, "Pong!", "Success")

		assert.Equal(t, expectedStatus, rec.Code)
		assert.Equal(t, expectedResp, rec.Body.String())
	})

	t.Run("should return 201 Created when use Created", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		expectedStatus := http.StatusCreated
		expectedResp := "{\"code\":\"0000\",\"status\":\"SUCCESS\",\"message\":\"\",\"data\":\"Created\"}\n"

		Created(ctx, "Created")

		assert.Equal(t, expectedStatus, rec.Code)
		assert.Equal(t, expectedResp, rec.Body.String())
	})

	t.Run("should return 201 Created with message when use CreatedWithMessage", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		expectedStatus := http.StatusCreated
		expectedResp := "{\"code\":\"0000\",\"status\":\"SUCCESS\",\"message\":\"Created successfully!\",\"data\":\"Created\"}\n"

		CreatedWithMessage(ctx, "Created successfully!", "Created")

		assert.Equal(t, expectedStatus, rec.Code)
		assert.Equal(t, expectedResp, rec.Body.String())
	})

	t.Run("should return 400 Bad Request when use Fail", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		expectedStatus := http.StatusBadRequest
		expectedResp := "{\"code\":\"4000\",\"status\":\"FAIL\",\"message\":\"Bad Request\"}\n"

		Fail(ctx, http.StatusBadRequest, "4000", "Bad Request")

		assert.Equal(t, expectedStatus, rec.Code)
		assert.Equal(t, expectedResp, rec.Body.String())
	})

	t.Run("should return 500 Internal Server Error when use FailWithError", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		expectedStatus := http.StatusInternalServerError
		expectedResp := "{\"code\":\"5000\",\"status\":\"FAIL\",\"message\":\"Internal Server Error\"}\n"

		err := InternalServerError("5000", "Internal Server Error", nil)
		FailWithError(ctx, err)

		assert.Equal(t, expectedStatus, rec.Code)
		assert.Equal(t, expectedResp, rec.Body.String())
	})
}
