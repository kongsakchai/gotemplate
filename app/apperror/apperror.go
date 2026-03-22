package apperror

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/errs"
	"github.com/labstack/echo/v5"
)

func ErrorHandler(ctx *echo.Context, err error) {
	c := ctx.Request().Context()

	var errMsg string
	var logs []slog.Attr
	var appErr app.Error

	if e, ok := err.(app.Error); ok {
		errMsg = "app error"
		logs = errs.Logs(e.Err)
		appErr = e
	} else {
		httpCode := http.StatusInternalServerError
		var sc echo.HTTPStatusCoder
		if errors.As(err, &sc) {
			if tmp := sc.StatusCode(); tmp != 0 {
				httpCode = tmp
			}
		}

		msg := ""
		switch e := sc.(type) {
		case json.Marshaler:
			b, _ := e.MarshalJSON()
			msg = string(b)
		case *echo.HTTPError:
			msg = e.Message
			if msg == "" {
				msg = http.StatusText(httpCode)
			}
		default:
			msg = http.StatusText(httpCode)
		}

		errMsg = "unexpected error"
		logs = errs.Logs(err)
		appErr = app.Error{HTTPCode: httpCode, Code: fmt.Sprintf("%d", httpCode), Message: msg}
	}

	ctx.Logger().LogAttrs(c, slog.LevelError, errMsg, logs...)
	err = app.Fail(ctx, appErr)

	if err != nil {
		ctx.Logger().ErrorContext(c, "error handler fail", "err", err.Error()) // rare case
	}
}
