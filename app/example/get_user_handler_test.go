package example

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/stretchr/testify/assert"
)

func TestHandlerGetUser(t *testing.T) {
	t.Run("should return bad request when name param is empty", func(t *testing.T) {
		storage := newMockStorager(t)
		h := NewHandler(storage)

		ctx, _ := app.NewMockContext(http.MethodGet, "/users", "")
		ctx.SetParamNames("name")
		ctx.SetParamValues("")

		err := h.GetUser(ctx)

		assert.Error(t, err)
		appErr, ok := err.(app.Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, appErr.HTTPCode)
		assert.Equal(t, "4001", appErr.Code)
	})

	t.Run("should return internal error when storage lookup fails", func(t *testing.T) {
		storage := newMockStorager(t)
		h := NewHandler(storage)

		storage.EXPECT().UserByName("john").Return(User{}, errors.New("db error"))

		ctx, _ := app.NewMockContext(http.MethodGet, "/users/john", "")
		ctx.SetParamNames("name")
		ctx.SetParamValues("john")

		err := h.GetUser(ctx)

		assert.Error(t, err)
		appErr, ok := err.(app.Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.HTTPCode)
		assert.Equal(t, "5001", appErr.Code)
	})

	t.Run("should return not found when user does not exist", func(t *testing.T) {
		storage := newMockStorager(t)
		h := NewHandler(storage)

		storage.EXPECT().UserByName("john").Return(User{}, nil)

		ctx, _ := app.NewMockContext(http.MethodGet, "/users/john", "")
		ctx.SetParamNames("name")
		ctx.SetParamValues("john")

		err := h.GetUser(ctx)

		assert.Error(t, err)
		appErr, ok := err.(app.Error)
		assert.True(t, ok)
		assert.Equal(t, "4002", appErr.Code)
	})

	t.Run("should return ok with user data when found", func(t *testing.T) {
		storage := newMockStorager(t)
		h := NewHandler(storage)

		expectedUser := User{FirstName: "john", LastName: "doe", Age: 30}
		storage.EXPECT().UserByName("john").Return(expectedUser, nil)

		ctx, rec := app.NewMockContext(http.MethodGet, "/users/john", "")
		ctx.SetParamNames("name")
		ctx.SetParamValues("john")

		err := h.GetUser(ctx)

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
		storage := newMockStorager(t)
		h := NewHandler(storage)

		storage.EXPECT().UserByName("john").Return(User{}, errors.New("db error"))

		_, err := h.getUser("john")

		assert.False(t, err.IsEmpty())
		assert.Equal(t, http.StatusInternalServerError, err.HTTPCode)
		assert.Equal(t, "5001", err.Code)
	})

	t.Run("should return not found when user not present", func(t *testing.T) {
		storage := newMockStorager(t)
		h := NewHandler(storage)

		storage.EXPECT().UserByName("john").Return(User{}, nil)

		_, err := h.getUser("john")

		assert.False(t, err.IsEmpty())
		assert.Equal(t, "4002", err.Code)
	})

	t.Run("should return user and empty error on success", func(t *testing.T) {
		storage := newMockStorager(t)
		h := NewHandler(storage)

		expectedUser := User{FirstName: "john", LastName: "doe", Age: 30}
		storage.EXPECT().UserByName("john").Return(expectedUser, nil)

		user, err := h.getUser("john")

		assert.True(t, err.IsEmpty())
		assert.Equal(t, expectedUser, user)
	})
}