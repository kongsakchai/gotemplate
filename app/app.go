package app

import (
	"context"
	"log/slog"
)

type Context interface {
	Query(key string) string
	Param(key string) string
	Bind(obj any) error

	JSON(code int, obj any) error
	OK(obj any) error
	OKWithMessage(message string, obj any) error
	Created(obj any) error
	CreatedWithMessage(message string, obj any) error
	NotFound(err Error) error
	InternalServer(err Error) error
	BadRequest(err Error) error

	Ctx() context.Context
	Get(key string) any
	Set(key string, value any)

	Logger() *slog.Logger
}

type Handler func(ctx Context) error

type Router interface {
	Start(addr string) error
	Shutdown(ctx context.Context) error
}

type Response struct {
	Code    string `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}
