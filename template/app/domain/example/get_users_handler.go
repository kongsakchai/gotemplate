package example

import (
	"github.com/kongsakchai/gotemplate/template/app"
	"github.com/labstack/echo/v5"
)

func (h *handler) GetUsers(ctx *echo.Context) error {
	users := h.storage.Users()
	return app.Ok(ctx, users)
}
