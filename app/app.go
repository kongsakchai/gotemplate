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
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Ok(ctx echo.Context, data any) error {
	return ctx.JSON(http.StatusOK, Response{
		Code:   SuccessCode,
		Status: SuccessStatus,
		Data:   data,
	})
}

func OkWithMessage(ctx echo.Context, data any, msg string) error {
	return ctx.JSON(http.StatusOK, Response{
		Code:    SuccessCode,
		Status:  SuccessStatus,
		Data:    data,
		Message: msg,
	})
}

func Created(ctx echo.Context, data any) error {
	return ctx.JSON(http.StatusCreated, Response{
		Code:   SuccessCode,
		Status: SuccessStatus,
		Data:   data,
	})
}

func CreatedWithMessage(ctx echo.Context, data any, msg string) error {
	return ctx.JSON(http.StatusCreated, Response{
		Code:    SuccessCode,
		Status:  SuccessStatus,
		Data:    data,
		Message: msg,
	})
}

func Fail(ctx echo.Context, err Error) error {
	if err.Err != nil {
		slog.Error("response error", "error", err.Err.Error())
	}

	return ctx.JSON(err.HTTPCode, Response{
		Code:    err.Code,
		Status:  ErrStatus,
		Data:    err.Data,
		Message: err.Message,
	})
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
