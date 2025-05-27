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
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Ok(ctx echo.Context, data any) error {
	return ctx.JSON(http.StatusOK, Response{
		Code:   SuccessCode,
		Status: SuccessStatus,
		Data:   data,
	})
}

func OkWithMessage(ctx echo.Context, msg string, data any) error {
	return ctx.JSON(http.StatusOK, Response{
		Code:    SuccessCode,
		Status:  SuccessStatus,
		Message: msg,
		Data:    data,
	})
}

func Created(ctx echo.Context, data any) error {
	return ctx.JSON(http.StatusCreated, Response{
		Code:   SuccessCode,
		Status: SuccessStatus,
		Data:   data,
	})
}

func CreatedWithMessage(ctx echo.Context, msg string, data any) error {
	return ctx.JSON(http.StatusCreated, Response{
		Code:    SuccessCode,
		Status:  SuccessStatus,
		Message: msg,
		Data:    data,
	})
}

func Fail(ctx echo.Context, statusCd int, code string, msg string) error {
	return ctx.JSON(statusCd, Response{
		Code:    code,
		Status:  FailureStatus,
		Message: msg,
	})
}

func FailWithError(ctx echo.Context, err Error) error {
	return ctx.JSON(err.StatusCd, Response{
		Code:    err.Code,
		Status:  FailureStatus,
		Message: err.Message,
	})
}
