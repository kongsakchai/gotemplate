package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/kongsakchai/gotemplate/errs"
	"github.com/labstack/echo/v4"
)

type App interface {
	Shutdown(ctx context.Context) error
	Start(addr string) error
}

type Response struct {
	Code    string `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Ok(ctx echo.Context, data any, msg ...string) error {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}
	return ctx.JSON(http.StatusOK, Response{
		Code:    SuccessCode,
		Success: true,
		Data:    data,
		Message: message,
	})
}

func Created(ctx echo.Context, data any, msg ...string) error {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}
	return ctx.JSON(http.StatusCreated, Response{
		Code:    SuccessCode,
		Success: true,
		Data:    data,
		Message: message,
	})
}

func Fail(ctx echo.Context, err Error) error {
	if err.Err != nil {
		slog.Error("response error", "error", err.Err.Error())
	}

	return ctx.JSON(err.HTTPCode, Response{
		Code:    err.Code,
		Success: false,
		Data:    err.Data,
		Message: err.Message,
	})
}

func ErrorHandler(err error, ctx echo.Context) {
	if appErr, ok := err.(Error); ok {
		err = Fail(ctx, appErr)
	} else {
		err = Fail(ctx, InternalServer(ErrInternalCode, ErrInternalMsg, err))
	}
	if err != nil {
		slog.Error("error handler fail", "err", err.Error())
	}
}

func logError(err error) {
	if err == nil {
		return
	}

	if wrapErr, ok := errs.As(err); ok {
		slog.Error("error", "error", wrapErr.UnwrapError(), "at", wrapErr.At())
	} else {
		slog.Error("error", "error", err.Error())
	}
}
