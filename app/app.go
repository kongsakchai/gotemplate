package app

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

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

	return ctx.JSON(err.StatusCd, Response{
		Code:    err.Code,
		Status:  FailureStatus,
		Message: err.Message,
	})
}

func FailWithData(ctx echo.Context, err Error, data any) error {
	if err.Error != nil {
		slog.Error("response error", "error", err.Error.Error())
	}

	return ctx.JSON(err.StatusCd, Response{
		Code:    err.Code,
		Status:  FailureStatus,
		Data:    data,
		Message: err.Message,
	})
}

// Query functions for extracting parameters from the request context

func QueryString(ctx echo.Context, key, defaultValue string) string {
	if value := ctx.QueryParam(key); value != "" {
		return value
	}
	return defaultValue
}

func QueryInt(ctx echo.Context, key string, defaultValue int64) int64 {
	if str := ctx.QueryParam(key); str != "" {
		if value, err := strconv.ParseInt(str, 10, 64); err == nil {
			return value
		}
	}
	return defaultValue
}

func QueryBool(ctx echo.Context, key string, defaultValue bool) bool {
	if str := ctx.QueryParam(key); str != "" {
		if value, err := strconv.ParseBool(str); err == nil {
			return value
		}
	}
	return defaultValue
}

func QueryFloat(ctx echo.Context, key string, defaultValue float64) float64 {
	if str := ctx.QueryParam(key); str != "" {
		if value, err := strconv.ParseFloat(str, 64); err == nil {
			return value
		}
	}
	return defaultValue
}

func ParamString(ctx echo.Context, key, defaultValue string) string {
	if value := ctx.Param(key); value != "" {
		return value
	}
	return defaultValue
}

func ParamInt(ctx echo.Context, key string, defaultValue int64) int64 {
	if str := ctx.Param(key); str != "" {
		if value, err := strconv.ParseInt(str, 10, 64); err == nil {
			return value
		}
	}
	return defaultValue
}
