package apperror

import (
	"fmt"
	"log/slog"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/errs"
	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, ctx echo.Context) {
	traceID, _ := ctx.Get(app.RefIDKey).(string)
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
		appErr = app.Error{HTTPCode: e.Code, Code: fmt.Sprintf("http %d", e.Code), Message: e.Message.(string)}
	default:
		msg = "service error"
		logs = errs.Logs(err)
		appErr = app.InternalError(app.ErrInternalCode, app.ErrInternalMsg, err)
	}

	slog.LogAttrs(c, slog.LevelError, msg, append(logs, slog.String(app.TraceIDKey, traceID))...)
	err = app.Fail(ctx, appErr)

	if err != nil {
		slog.ErrorContext(c, "error handler fail", "err", err.Error())
	}
}
