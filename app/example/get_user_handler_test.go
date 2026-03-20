package example

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
)

func TestHandlerGetUser(t *testing.T) {
	t.Run("should return bad request when name param is empty", func(t *testing.T) {
		// arrange
		storage := newMockStorager(t)
		h := NewHandler(storage)

		ctx := echotest.ContextConfig{
			Headers: http.Header{
				echo.HeaderContentType: []string{echo.MIMEApplicationJSON},
			},
			PathValues: echo.PathValues{
				{Name: "name", Value: ""},
			},
		}.ToContext(t)

		// act
		err := h.GetUser(ctx)

		// assert
		assert.Error(t, err)
		appErr, ok := err.(app.Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, appErr.HTTPCode)
		assert.Equal(t, "4001", appErr.Code)
	})

	t.Run("should return internal error when storage lookup fails", func(t *testing.T) {
		// arrange
		storage := newMockStorager(t)
		h := NewHandler(storage)

		storage.EXPECT().UserByName("john").Return(User{}, errors.New("db error"))

		ctx := echotest.ContextConfig{
			Headers: http.Header{
				echo.HeaderContentType: []string{echo.MIMEApplicationJSON},
			},
			PathValues: echo.PathValues{
				{Name: "name", Value: "john"},
			},
		}.ToContext(t)

		// act
		err := h.GetUser(ctx)

		// assert
		assert.Error(t, err)
		appErr, ok := err.(app.Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.HTTPCode)
		assert.Equal(t, "5001", appErr.Code)
	})

	t.Run("should return not found when user does not exist", func(t *testing.T) {
		// arrange
		storage := newMockStorager(t)
		h := NewHandler(storage)

		storage.EXPECT().UserByName("john").Return(User{}, nil)

		ctx := echotest.ContextConfig{
			Headers: http.Header{
				echo.HeaderContentType: []string{echo.MIMEApplicationJSON},
			},
			PathValues: echo.PathValues{
				{Name: "name", Value: "john"},
			},
		}.ToContext(t)

		// act
		err := h.GetUser(ctx)

		// assert
		assert.Error(t, err)
		appErr, ok := err.(app.Error)
		assert.True(t, ok)
		assert.Equal(t, "4002", appErr.Code)
	})

	t.Run("should return ok with user data when found", func(t *testing.T) {
		// arrange
		storage := newMockStorager(t)
		h := NewHandler(storage)

		expectedUser := User{FirstName: "john", LastName: "doe", Age: 30}
		storage.EXPECT().UserByName("john").Return(expectedUser, nil)

		ctx, rec := echotest.ContextConfig{
			Headers: http.Header{
				echo.HeaderContentType: []string{echo.MIMEApplicationJSON},
			},
			PathValues: echo.PathValues{
				{Name: "name", Value: "john"},
			},
		}.ToContextRecorder(t)

		// act
		err := h.GetUser(ctx)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp app.Response
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.True(t, resp.Success)
		assert.Equal(t, app.SuccessCode, resp.Code)
	})
}

func TestHandlerGetUserFunction(t *testing.T) {
	t.Run("should return internal error when user lookup fails", func(t *testing.T) {
		// arrange
		storage := newMockStorager(t)
		h := NewHandler(storage)

		storage.EXPECT().UserByName("john").Return(User{}, errors.New("db error"))

		// act
		_, err := h.getUser("john")

		// assert
		assert.False(t, err.IsEmpty())
		assert.Equal(t, http.StatusInternalServerError, err.HTTPCode)
		assert.Equal(t, "5001", err.Code)
	})

	t.Run("should return not found when user not present", func(t *testing.T) {
		// arrange
		storage := newMockStorager(t)
		h := NewHandler(storage)

		storage.EXPECT().UserByName("john").Return(User{}, nil)

		// act
		_, err := h.getUser("john")

		// assert
		assert.False(t, err.IsEmpty())
		assert.Equal(t, "4002", err.Code)
	})

	t.Run("should return user and empty error on success", func(t *testing.T) {
		// arrange
		storage := newMockStorager(t)
		h := NewHandler(storage)

		expectedUser := User{FirstName: "john", LastName: "doe", Age: 30}
		storage.EXPECT().UserByName("john").Return(expectedUser, nil)

		// act
		user, err := h.getUser("john")

		// assert
		assert.True(t, err.IsEmpty())
		assert.Equal(t, expectedUser, user)
	})
}
