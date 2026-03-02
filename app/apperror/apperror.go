package apperror

import (
	"log/slog"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/errs"
	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, ctx echo.Context) {
	traceID, _ := ctx.Get(app.RefIDKey).(string)

	c := ctx.Request().Context()
	if appErr, ok := err.(app.Error); ok {
		slog.LogAttrs(c, slog.LevelError, "app error", append(errs.Logs(appErr.Err), slog.String(app.TraceIDKey, traceID))...)
		err = app.Fail(ctx, appErr)
	} else {
		slog.LogAttrs(c, slog.LevelError, "app error", append(errs.Logs(err), slog.String(app.TraceIDKey, traceID))...)
		err = app.Fail(ctx, app.InternalError(app.ErrInternalCode, app.ErrInternalMsg, err))
	}
	if err != nil {
		slog.ErrorContext(c, "error handler fail", "err", err.Error())
	}
}
