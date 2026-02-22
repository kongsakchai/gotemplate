package app

import (
	"context"
	"net/http"

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
	return ctx.JSON(err.HTTPCode, Response{
		Code:    err.Code,
		Success: false,
		Data:    err.Data,
		Message: err.Message,
	})
}
