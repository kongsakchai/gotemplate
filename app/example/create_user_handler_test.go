package example

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/kongsakchai/gotemplate/app"
	appValidator "github.com/kongsakchai/gotemplate/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandlerCreateUser(t *testing.T) {
	t.Run("should return bad request when body is invalid", func(t *testing.T) {
		storage := newMockStorager(t)
		h := NewHandler(storage)

		ctx, _ := app.NewMockContext(http.MethodPost, "/users", "{invalid-json")
		ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx.Echo().Validator = appValidator.NewReqValidator()

		err := h.CreateUser(ctx)

		assert.Error(t, err)
		appErr, ok := err.(app.Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, appErr.HTTPCode)
		assert.Equal(t, "4001", appErr.Code)
	})

	t.Run("should return bad request when validation fails", func(t *testing.T) {
		storage := newMockStorager(t)
		h := NewHandler(storage)

		ctx, _ := app.NewMockContext(http.MethodPost, "/users", `{}`)
		ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx.Echo().Validator = appValidator.NewReqValidator()

		err := h.CreateUser(ctx)

		assert.Error(t, err)
		appErr, ok := err.(app.Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, appErr.HTTPCode)
		assert.Equal(t, "4002", appErr.Code)
	})

	t.Run("should return conflict when user already exists", func(t *testing.T) {
		storage := newMockStorager(t)
		h := NewHandler(storage)

		storage.EXPECT().UserByName("john").Return(User{FirstName: "john"}, nil)

		ctx, _ := app.NewMockContext(http.MethodPost, "/users", `{"firstName":"john","lastName":"doe","age":30}`)
		ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx.Echo().Validator = appValidator.NewReqValidator()

		err := h.CreateUser(ctx)

		assert.Error(t, err)
		appErr, ok := err.(app.Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusConflict, appErr.HTTPCode)
		assert.Equal(t, "4003", appErr.Code)
	})

	t.Run("should return ok when create user successful", func(t *testing.T) {
		storage := newMockStorager(t)
		h := NewHandler(storage)

		storage.EXPECT().UserByName("john").Return(User{}, nil)
		storage.EXPECT().CreateUser(User{FirstName: "john", LastName: "doe", Age: 30}).Return(nil)

		ctx, rec := app.NewMockContext(http.MethodPost, "/users", `{"firstName":"john","lastName":"doe","age":30}`)
		ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx.Echo().Validator = appValidator.NewReqValidator()

		err := h.CreateUser(ctx)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp app.Response
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.True(t, resp.Success)
		assert.Equal(t, app.SuccessCode, resp.Code)
	})
}

func TestHandlerCreateUserFunction(t *testing.T) {
	t.Run("should return internal error when lookup fails", func(t *testing.T) {
		storage := newMockStorager(t)
		h := NewHandler(storage)

		storage.EXPECT().UserByName("john").Return(User{}, errors.New("db error"))

		err := h.createUser(CreateUserRequest{FirstName: "john", LastName: "doe", Age: 30})

		assert.False(t, err.IsEmpty())
		assert.Equal(t, http.StatusInternalServerError, err.HTTPCode)
		assert.Equal(t, "5001", err.Code)
	})

	t.Run("should return internal error when create fails", func(t *testing.T) {
		storage := newMockStorager(t)
		h := NewHandler(storage)

		storage.EXPECT().UserByName("john").Return(User{}, nil)
		storage.EXPECT().CreateUser(User{FirstName: "john", LastName: "doe", Age: 30}).Return(errors.New("insert error"))

		err := h.createUser(CreateUserRequest{FirstName: "john", LastName: "doe", Age: 30})

		assert.False(t, err.IsEmpty())
		assert.Equal(t, http.StatusInternalServerError, err.HTTPCode)
		assert.Equal(t, "5002", err.Code)
	})

	t.Run("should return empty error on success", func(t *testing.T) {
		storage := newMockStorager(t)
		h := NewHandler(storage)

		storage.EXPECT().UserByName("john").Return(User{}, nil)
		storage.EXPECT().CreateUser(User{FirstName: "john", LastName: "doe", Age: 30}).Return(nil)

		err := h.createUser(CreateUserRequest{FirstName: "john", LastName: "doe", Age: 30})

		assert.True(t, err.IsEmpty())
	})
}