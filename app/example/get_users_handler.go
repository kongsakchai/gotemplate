package example

import (
	"github.com/kongsakchai/gotemplate/app"
	"github.com/labstack/echo/v4"
)

func (h *handler) GetUserByName(ctx echo.Context) error {
	name := ctx.Param("name")
	user, err := h.storage.UserByName(name)
	if !err.IsEmpty() {
		return err
	}
	return app.Ok(ctx, user)
}

func (h *handler) GetUsers(ctx echo.Context) error {
	users, err := h.storage.Users()
	if !err.IsEmpty() {
		return err
	}
	return app.Ok(ctx, users)
}
