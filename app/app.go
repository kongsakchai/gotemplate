package app

import (
	"context"
	"log/slog"
	"net/http"

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
	if err.Error != nil {
		slog.Error("response error", "error", err.Error.Error())
	}

	return ctx.JSON(err.HTTPCode, Response{
		Code:    err.Code,
		Status:  ErrStatus,
		Message: err.Message,
	})
}

func FailWithData(ctx echo.Context, err Error, data any) error {
	logError(err.Error)

	return ctx.JSON(err.HTTPCode, Response{
		Code:    err.Code,
		Status:  ErrStatus,
		Data:    data,
		Message: err.Message,
	})
}

type WrapError interface {
	Error() string
	AtError() string
	RawError() string
}

func logError(err error) {
	if err == nil {
		return
	}

	if wrapErr, ok := err.(WrapError); ok {
		slog.Error("error", "error", wrapErr.RawError(), "at", wrapErr.AtError())
	} else {
		slog.Error("error", "error", err.Error())
	}
}
