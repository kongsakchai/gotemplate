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
		expectedResp := "{\"code\":\"0000\",\"status\":\"SUCCESS\",\"data\":\"Success\"}\n"

		Ok(ctx, "Success")

		assert.Equal(t, expectedStatus, rec.Code)
		assert.JSONEq(t, expectedResp, rec.Body.String())
	})

	t.Run("should return 200 OK with message when use OkWithMessage", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		expectedStatus := http.StatusOK
		expectedResp := "{\"code\":\"0000\",\"status\":\"SUCCESS\",\"message\":\"Pong!\",\"data\":\"Success\"}\n"

		OkWithMessage(ctx, "Success", "Pong!")

		assert.Equal(t, expectedStatus, rec.Code)
		assert.JSONEq(t, expectedResp, rec.Body.String())
	})

	t.Run("should return 201 Created when use Created", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		expectedStatus := http.StatusCreated
		expectedResp := "{\"code\":\"0000\",\"status\":\"SUCCESS\",\"data\":\"Created\"}\n"

		Created(ctx, "Created")

		assert.Equal(t, expectedStatus, rec.Code)
		assert.JSONEq(t, expectedResp, rec.Body.String())
	})

	t.Run("should return 201 Created with message when use CreatedWithMessage", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		expectedStatus := http.StatusCreated
		expectedResp := "{\"code\":\"0000\",\"status\":\"SUCCESS\",\"message\":\"Created successfully!\",\"data\":\"Created\"}\n"

		CreatedWithMessage(ctx, "Created", "Created successfully!")

		assert.Equal(t, expectedStatus, rec.Code)
		assert.JSONEq(t, expectedResp, rec.Body.String())
	})

	t.Run("should return 500 Internal Server Error when use FailWithError", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		expectedStatus := http.StatusInternalServerError
		expectedResp := "{\"code\":\"5000\",\"status\":\"FAIL\",\"message\":\"Internal Server Error\"}\n"

		err := InternalServerError("5000", "Internal Server Error", errors.New("unexpected error"))
		Fail(ctx, err)

		assert.Equal(t, expectedStatus, rec.Code)
		assert.JSONEq(t, expectedResp, rec.Body.String())
	})

	t.Run("should return 400 Bad Request with error message when use FailWithData", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		expectedStatus := http.StatusBadRequest
		expectedResp := "{\"code\":\"4000\",\"status\":\"FAIL\",\"message\":\"Bad Request\",\"data\":\"Invalid data\"}\n"

		err := BadRequestError("4000", "Bad Request", errors.New("invalid input"))
		FailWithData(ctx, err, "Invalid data")

		assert.Equal(t, expectedStatus, rec.Code)
		assert.JSONEq(t, expectedResp, rec.Body.String())
	})
}

func TestMakeResponse(t *testing.T) {
	t.Run("should create response display with title and message", func(t *testing.T) {
		display := []string{"Title", "Message"}
		respDisplay := makeResponseDisplay(display)

		assert.Equal(t, "Title", respDisplay.Title)
		assert.Equal(t, "Message", respDisplay.Message)
	})

	t.Run("should create response display with only title", func(t *testing.T) {
		display := []string{"Only Title"}
		respDisplay := makeResponseDisplay(display)

		assert.Equal(t, "Only Title", respDisplay.Title)
		assert.Equal(t, "", respDisplay.Message)
	})

	t.Run("should create empty response display", func(t *testing.T) {
		display := []string{}
		respDisplay := makeResponseDisplay(display)

		assert.Equal(t, "", respDisplay.Title)
		assert.Equal(t, "", respDisplay.Message)
	})
}
