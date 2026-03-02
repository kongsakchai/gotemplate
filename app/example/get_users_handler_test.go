package example

import (
	"testing"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByName(t *testing.T) {
	t.Run("should return name when service success", func(t *testing.T) {
		service := newMockStorager(t)
		service.EXPECT().UserByName("john").Return(User{Name: "john"}, app.Error{})

		h := NewHandler(service)

		ctx, rec := app.NewMockContext("GET", "/john", "")
		ctx.SetParamNames("name")
		ctx.SetParamValues("john")

		err := h.GetUserByName(ctx)

		assert.NoError(t, err)
		assert.JSONEq(t, `{"code":"0000","data":{"name":"john","age":0},"success":true}`, rec.Body.String())
	})

	t.Run("should return error when service returns error", func(t *testing.T) {
		service := newMockStorager(t)
		service.EXPECT().UserByName("john").Return(User{}, app.NotFound("4001", "user not found", nil))

		h := NewHandler(service)

		ctx, _ := app.NewMockContext("GET", "/john", "")
		ctx.SetParamNames("name")
		ctx.SetParamValues("john")

		err := h.GetUserByName(ctx)

		assert.Error(t, err)
		assert.Equal(t, "http_code=400 code=4001 msg=\"user not found\" data=<nil> err=<nil>", err.Error())
	})
}

func TestGetUsers(t *testing.T) {
	t.Run("should return users when service success", func(t *testing.T) {
		service := newMockStorager(t)
		service.EXPECT().Users().Return([]User{{Name: "john"}}, app.Error{})

		h := NewHandler(service)

		ctx, rec := app.NewMockContext("GET", "/users", "")

		err := h.GetUsers(ctx)

		assert.NoError(t, err)
		assert.JSONEq(t, `{"code":"0000","data":[{"name":"john","age":0}],"success":true}`, rec.Body.String())
	})

	t.Run("should return error when service returns error", func(t *testing.T) {
		service := newMockStorager(t)
		service.EXPECT().Users().Return(nil, app.InternalError("5001", "internal server error", nil))

		h := NewHandler(service)

		ctx, _ := app.NewMockContext("GET", "/users", "")

		err := h.GetUsers(ctx)

		assert.Error(t, err)
		assert.Equal(t, "http_code=500 code=5001 msg=\"internal server error\" data=<nil> err=<nil>", err.Error())
	})
}
