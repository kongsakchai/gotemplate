package app

import (
	"context"
	"log/slog"
)

type Context interface {
	Next(ctx Context) error

	Query(key string) string
	Param(key string) string
	Bind(obj any) error

	JSON(code int, obj any) error
	OK(obj any) error
	Created(obj any) error
	NotFound(err Error) error
	InternalServer(err Error) error
	BadRequest(err Error) error

	Ctx() context.Context
	Get(key string) any
	Set(key string, value any)

	CtxLogger() *slog.Logger
	SetCtxLogger(logger *slog.Logger)
}

type Handler func(ctx Context) error

type Router interface {
	GET(path string, handler Handler)
	POST(path string, handler Handler)
	PUT(path string, handler Handler)
	DELETE(path string, handler Handler)
	Use(middleware ...Handler)
	Shutdown(ctx context.Context) error
	Start(addr string) error
}

type Response struct {
	Code    string `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}
