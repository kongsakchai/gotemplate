package app

import (
	"errors"
	"net/http"
	"testing"

	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
)

type mockRequestContext struct {
	bindFn     func(any) error
	validateFn func(any) error
}

func (m *mockRequestContext) Bind(i any) error {
	if m.bindFn != nil {
		return m.bindFn(i)
	}
	return nil
}

func (m *mockRequestContext) Validate(i any) error {
	if m.validateFn != nil {
		return m.validateFn(i)
	}
	return nil
}

func TestRequest(t *testing.T) {
	t.Run("should return bad request when Bind fails", func(t *testing.T) {
		ctx := &mockRequestContext{
			bindFn: func(i any) error {
				return errors.New("bind error")
			},
		}

		err := Request(ctx, &struct{}{})
		assert.Error(t, err)

		appErr, ok := err.(Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, appErr.HTTPCode)
		assert.Equal(t, BadRequestCode, appErr.Code)
	})

	t.Run("should return invalid request when Validate fails", func(t *testing.T) {
		ctx := &mockRequestContext{
			bindFn: func(i any) error {
				return nil
			},
			validateFn: func(i any) error {
				return errors.New("validate error")
			},
		}

		err := Request(ctx, &struct{}{})
		assert.Error(t, err)

		appErr, ok := err.(Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, appErr.HTTPCode)
		assert.Equal(t, InValidCode, appErr.Code)
	})

	t.Run("should return nil when Bind and Validate succeed", func(t *testing.T) {
		ctx := &mockRequestContext{
			bindFn: func(i any) error {
				return nil
			},
			validateFn: func(i any) error {
				return nil
			},
		}

		err := Request(ctx, &struct{}{})
		assert.NoError(t, err)
	})
}

func TestAppResponse(t *testing.T) {
	t.Run("should return 200 OK when use Ok", func(t *testing.T) {
		// arrange
		ctx, rec := echotest.ContextConfig{}.ToContextRecorder(t)

		expectedStatus := http.StatusOK
		expectedResp := "{\"code\":\"0000\",\"success\":true,\"data\":\"Success\",\"message\":\"Ping!\"}\n"

		// act
		Ok(ctx, "Success", "Ping!")

		// assert
		assert.Equal(t, expectedStatus, rec.Code)
		assert.JSONEq(t, expectedResp, rec.Body.String())
	})

	t.Run("should return 201 Created when use Created", func(t *testing.T) {
		// arrange
		ctx, rec := echotest.ContextConfig{}.ToContextRecorder(t)

		expectedStatus := http.StatusCreated
		expectedResp := "{\"code\":\"0000\",\"success\":true,\"data\":\"Created\",\"message\":\"Ping!\"}\n"

		// act
		Created(ctx, "Created", "Ping!")

		// assert
		assert.Equal(t, expectedStatus, rec.Code)
		assert.JSONEq(t, expectedResp, rec.Body.String())
	})

	t.Run("should return 500 Internal Server Error when use FailWithError", func(t *testing.T) {
		// arrange
		ctx, rec := echotest.ContextConfig{}.ToContextRecorder(t)

		expectedStatus := http.StatusInternalServerError
		expectedResp := "{\"code\":\"5000\",\"success\":false,\"message\":\"Internal Server Error\"}\n"

		// act
		err := InternalError("5000", "Internal Server Error", errors.New("unexpected error"))
		Fail(ctx, err)

		// assert
		assert.Equal(t, expectedStatus, rec.Code)
		assert.JSONEq(t, expectedResp, rec.Body.String())
	})
}
