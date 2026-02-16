package app

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	t.Run("should return 400 Bad Request when use BadRequest", func(t *testing.T) {
		expectedError := Error{
			HTTPCode: http.StatusBadRequest,
			Code:     "4000",
			Message:  "Bad Request",
			Data:     "hello",
			Err:      nil,
		}

		err := BadRequest("4000", "Bad Request", nil, "hello")

		assert.Equal(t, expectedError, err)
	})

	t.Run("should return 404 Not Found when use NotFound", func(t *testing.T) {
		expectedError := Error{
			HTTPCode: http.StatusOK,
			Code:     "4040",
			Message:  "Not Found",
			Err:      nil,
		}

		err := NotFound("4040", "Not Found", nil)

		assert.Equal(t, expectedError, err)
	})

	t.Run("should return 500 Internal Server Error when use InternalServerError", func(t *testing.T) {
		expectedError := Error{
			HTTPCode: http.StatusInternalServerError,
			Code:     "5000",
			Message:  "Internal Server Error",
			Err:      nil,
		}

		err := InternalServer("5000", "Internal Server Error", nil)

		assert.Equal(t, expectedError, err)
	})

	t.Run("should return 401 Unauthorized when use UnauthorizedError", func(t *testing.T) {
		expectedError := Error{
			HTTPCode: http.StatusUnauthorized,
			Code:     "4010",
			Message:  "Unauthorized",
			Err:      nil,
		}

		err := Unauthorized("4010", "Unauthorized", nil)

		assert.Equal(t, expectedError, err)
	})

	t.Run("should return 403 Forbidden when use ForbiddenError", func(t *testing.T) {
		expectedError := Error{
			HTTPCode: http.StatusForbidden,
			Code:     "4030",
			Message:  "Forbidden",
			Err:      nil,
		}

		err := Forbidden("4030", "Forbidden", nil)

		assert.Equal(t, expectedError, err)
	})

	t.Run("should return 409 Conflict when use ConflictError", func(t *testing.T) {
		expectedError := Error{
			HTTPCode: http.StatusConflict,
			Code:     "4090",
			Message:  "Conflict",
			Err:      nil,
		}

		err := Conflict("4090", "Conflict", nil)

		assert.Equal(t, expectedError, err)
	})

	t.Run("should return true when use IsEmpty", func(t *testing.T) {
		err := Error{}
		assert.True(t, err.IsEmpty())

		err = BadRequest("4000", "Bad Request", nil)
		assert.False(t, err.IsEmpty())
	})

	t.Run("should return error when call Error()", func(t *testing.T) {
		err := Error{}
		assert.True(t, err.IsEmpty())

		err = BadRequest("4000", "Bad Request", nil)
		assert.Equal(t, "http_code=400 code=4000 msg=\"Bad Request\" data=<nil> err=<nil>", err.Error())
	})
}

func TestErrorData(t *testing.T) {
	t.Run("should return nil when get nil data", func(t *testing.T) {
		data := errorData(nil)
		assert.Nil(t, data)
	})

	t.Run("should return data when get one data", func(t *testing.T) {
		data := errorData([]any{"hello"})
		assert.Equal(t, "hello", data)
	})

	t.Run("should return array when get multiple data", func(t *testing.T) {
		data := errorData([]any{"hello", "world"})
		assert.Equal(t, []any{"hello", "world"}, data)
	})
}
