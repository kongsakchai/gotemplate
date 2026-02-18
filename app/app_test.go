package app

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kongsakchai/gotemplate/errs"
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

		err := InternalServer("5000", "Internal Server Error", errors.New("unexpected error"))
		Fail(ctx, err)

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
			err:      errs.New("error"),
			expected: "level=ERROR msg=error error=error at=",
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

type mockFailContext struct {
	echo.Context
}

func (m mockFailContext) JSON(code int, i any) error {
	return errors.New("JSON ERROR")
}

func TestErrorHandler(t *testing.T) {
	t.Run("should return bad request when error is app error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		appErr := BadRequest("4000", "error", errors.New("error"))
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

	t.Run("should no response when send response fail", func(t *testing.T) {
		ctx := &mockFailContext{}

		appErr := BadRequest("4000", "error", errors.New("error"))
		ErrorHandler(appErr, ctx)
	})
}
