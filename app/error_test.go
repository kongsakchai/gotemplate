package app

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	t.Run("should return 400 Bad Request when use BadRequest", func(t *testing.T) {
		expectedError := Error{
			StatusCd: http.StatusBadRequest,
			Code:     "4000",
			Message:  "Bad Request",
			Error:    nil,
		}

		err := BadRequestError("4000", "Bad Request", nil)

		assert.Equal(t, expectedError, err)
	})

	t.Run("should return 404 Not Found when use NotFound", func(t *testing.T) {
		expectedError := Error{
			StatusCd: http.StatusOK,
			Code:     "4040",
			Message:  "Not Found",
			Error:    nil,
		}

		err := NotFoundError("4040", "Not Found", nil)

		assert.Equal(t, expectedError, err)
	})

	t.Run("should return 500 Internal Server Error when use InternalServerError", func(t *testing.T) {
		expectedError := Error{
			StatusCd: http.StatusInternalServerError,
			Code:     "5000",
			Message:  "Internal Server Error",
			Error:    nil,
		}

		err := InternalServerError("5000", "Internal Server Error", nil)

		assert.Equal(t, expectedError, err)
	})

	t.Run("should return 401 Unauthorized when use UnauthorizedError", func(t *testing.T) {
		expectedError := Error{
			StatusCd: http.StatusUnauthorized,
			Code:     "4010",
			Message:  "Unauthorized",
			Error:    nil,
		}

		err := UnauthorizedError("4010", "Unauthorized", nil)

		assert.Equal(t, expectedError, err)
	})

	t.Run("should return 403 Forbidden when use ForbiddenError", func(t *testing.T) {
		expectedError := Error{
			StatusCd: http.StatusForbidden,
			Code:     "4030",
			Message:  "Forbidden",
			Error:    nil,
		}

		err := ForbiddenError("4030", "Forbidden", nil)

		assert.Equal(t, expectedError, err)
	})

	t.Run("should return 409 Conflict when use ConflictError", func(t *testing.T) {
		expectedError := Error{
			StatusCd: http.StatusConflict,
			Code:     "4090",
			Message:  "Conflict",
			Error:    nil,
		}

		err := ConflictError("4090", "Conflict", nil)

		assert.Equal(t, expectedError, err)
	})

	t.Run("should return true when use IsEmpty", func(t *testing.T) {
		err := Error{}
		assert.True(t, err.IsEmpty())

		err = BadRequestError("4000", "Bad Request", nil)
		assert.False(t, err.IsEmpty())
	})
}
