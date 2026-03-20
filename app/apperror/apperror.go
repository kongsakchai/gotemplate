package apperror

import (
	"fmt"
	"log/slog"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/errs"
	"github.com/labstack/echo/v5"
)

func ErrorHandler(ctx *echo.Context, err error) {
	c := ctx.Request().Context()

	var msg string
	var logs []slog.Attr
	var appErr app.Error

	switch e := err.(type) {
	case app.Error:
		msg = "app error"
		logs = errs.Logs(e.Err)
		appErr = e
	case *echo.HTTPError:
		msg = "http error"
		logs = errs.Logs(e)
		appErr = app.Error{HTTPCode: e.Code, Code: fmt.Sprintf("http %d", e.Code), Message: e.Message}
	default:
		msg = "service error"
		logs = errs.Logs(err)
		appErr = app.InternalError(app.ErrInternalCode, app.ErrInternalMsg, err)
	}

	ctx.Logger().LogAttrs(c, slog.LevelError, msg, logs...)
	err = app.Fail(ctx, appErr)

	if err != nil {
		ctx.Logger().ErrorContext(c, "error handler fail", "err", err.Error()) // rare case
	}
}
