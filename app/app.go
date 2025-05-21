package app

import (
	"context"
	"log/slog"
)

type Validator interface {
	Validate(obj any) error
}

type Context interface {
	Validator

	Next() error
	Query(key string) string
	Param(key string) string
	Bind(obj any) error

	JSON(code int, obj any) error
	OK(obj any) error
	OKWithMessage(message string, obj any) error
	Created(obj any) error
	CreatedWithMessage(message string, obj any) error
	Error(err *Error) error

	Ctx() context.Context
	Get(key string) any
	Set(key string, value any)

	Logger() *slog.Logger
}

type Handler func(ctx Context) error

type App interface {
	Start(addr string) error
	Shutdown(ctx context.Context) error
	Use(middlewares ...Handler)
	GET(path string, handler ...Handler)
	POST(path string, handler ...Handler)
	PUT(path string, handler ...Handler)
	DELETE(path string, handler ...Handler)
	PATCH(path string, handler ...Handler)
	Group(path string, handlers ...Handler) AppGroup
}

type AppGroup interface {
	Use(middlewares ...Handler)
	GET(path string, handler ...Handler)
	POST(path string, handler ...Handler)
	PUT(path string, handler ...Handler)
	DELETE(path string, handler ...Handler)
	PATCH(path string, handler ...Handler)
	Group(path string, handlers ...Handler) AppGroup
}

type Response struct {
	Code    string `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}
