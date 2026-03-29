package app

import (
	"errors"
	"net/http"
	"testing"

	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
)

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
