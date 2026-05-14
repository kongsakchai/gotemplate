package example

import (
	"github.com/kongsakchai/gotemplate/template/app"
	"github.com/labstack/echo/v5"
)

func (h *handler) GetUser(ctx *echo.Context) error {
	name := ctx.Param("name")
	if name == "" {
		return app.BadRequest("4001", "name parameter is required", nil)
	}

	user, err := h.getUser(name)
	if !err.IsEmpty() {
		return err
	}
	return app.Ok(ctx, user)
}

func (h *handler) getUser(name string) (User, app.Error) {
	user, err := h.storage.UserByName(name)
	if err != nil {
		return User{}, app.InternalError("5001", "failed to get user by name", err)
	}
	if user.FirstName == "" {
		return User{}, app.NotFound("4002", "user not found", nil)
	}
	return user, app.Error{}
}
