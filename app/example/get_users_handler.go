package example

import (
	"github.com/kongsakchai/gotemplate/app"
	"github.com/labstack/echo/v4"
)

func (h *handler) GetUsers(ctx echo.Context) error {
	users, err := h.storage.Users()
	if err != nil {
		return app.InternalError("5001", "failed to get users", err)
	}
	return app.Ok(ctx, users)
}
