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
	Code    string          `json:"code"`
	Status  string          `json:"status"`
	Message string          `json:"message,omitempty"`
	Data    any             `json:"data,omitempty"`
	Display ResponseDisplay `json:"display,omitzero"`
}

type ResponseDisplay struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func makeResponseDisplay(display []string) ResponseDisplay {
	var title, message string
	if len(display) > 0 {
		title = display[0]
	}
	if len(display) > 1 {
		message = display[1]
	}
	return ResponseDisplay{
		Title:   title,
		Message: message,
	}
}

func Ok(ctx echo.Context, data any, display ...string) error {
	return ctx.JSON(http.StatusOK, Response{
		Code:    SuccessCode,
		Status:  SuccessStatus,
		Data:    data,
		Display: makeResponseDisplay(display),
	})
}

func OkWithMessage(ctx echo.Context, data any, msg string, display ...string) error {
	return ctx.JSON(http.StatusOK, Response{
		Code:    SuccessCode,
		Status:  SuccessStatus,
		Data:    data,
		Message: msg,
		Display: makeResponseDisplay(display),
	})
}

func Created(ctx echo.Context, data any, display ...string) error {
	return ctx.JSON(http.StatusCreated, Response{
		Code:    SuccessCode,
		Status:  SuccessStatus,
		Data:    data,
		Display: makeResponseDisplay(display),
	})
}

func CreatedWithMessage(ctx echo.Context, data any, msg string, display ...string) error {
	return ctx.JSON(http.StatusCreated, Response{
		Code:    SuccessCode,
		Status:  SuccessStatus,
		Data:    data,
		Message: msg,
		Display: makeResponseDisplay(display),
	})
}

func Fail(ctx echo.Context, err Error, display ...string) error {
	if err.Error != nil {
		slog.Error("response error", "error", err.Error.Error())
	}

	return ctx.JSON(err.StatusCd, Response{
		Code:    err.Code,
		Status:  FailureStatus,
		Message: err.Message,
		Display: makeResponseDisplay(display),
	})
}

func FailWithData(ctx echo.Context, err Error, data any, display ...string) error {
	if err.Error != nil {
		slog.Error("response error", "error", err.Error.Error())
	}

	return ctx.JSON(err.StatusCd, Response{
		Code:    err.Code,
		Status:  FailureStatus,
		Data:    data,
		Message: err.Message,
		Display: makeResponseDisplay(display),
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
