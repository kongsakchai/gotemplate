package app

import (
	"bytes"
	"errors"
	"log/slog"
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
		expectedResp := "{\"code\":\"5000\",\"status\":\"ERROR\",\"message\":\"Internal Server Error\"}\n"

		err := InternalServer("5000", "Internal Server Error", errors.New("unexpected error"))
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
		expectedResp := "{\"code\":\"4000\",\"status\":\"ERROR\",\"message\":\"Bad Request\",\"data\":\"Invalid data\"}\n"

		err := BadRequest("4000", "Bad Request", errors.New("invalid input"))
		FailWithData(ctx, err, "Invalid data")

		assert.Equal(t, expectedStatus, rec.Code)
		assert.JSONEq(t, expectedResp, rec.Body.String())
	})
}

type mockError struct {
	Message string
}

func (e *mockError) Error() string {
	return e.Message
}

func (e *mockError) At() string {
	return "mock error location"
}

func (e *mockError) UnwrapError() string {
	return e.Message
}

func TestLogError(t *testing.T) {
	type testcase struct {
		title    string
		err      error
		expected string
	}

	testCases := []testcase{
		{
			title:    "should log error with message and at location",
			err:      errors.New("test error"),
			expected: "test error",
		},
		{
			title:    "should log error with nil error",
			err:      nil,
			expected: "",
		},
		{
			title:    "should log custom error with message and at location",
			err:      &mockError{Message: "custom error"},
			expected: "error=\"custom error\" at=\"mock error location\"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			df := slog.Default()
			defer func() {
				slog.SetDefault(df)
			}()

			wt := bytes.NewBuffer([]byte{})
			slog.SetDefault(slog.New(slog.NewTextHandler(wt, nil)))

			logError(tc.err)

			if tc.expected == "" {
				assert.Empty(t, wt.String())
			} else {
				assert.Contains(t, wt.String(), tc.expected)
			}
		})
	}
}
