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
	Created(obj any) error
	NotFound(err Error) error
	InternalServer(err Error) error
	BadRequest(err Error) error

	Ctx() context.Context
	Get(key string) any
	Set(key string, value any)

	Logger() *slog.Logger
	SetLogger(logger *slog.Logger)
}

type Handler func(ctx Context) error

type Middleware func(next Handler) Handler

type Router interface {
	Origin() any
	GET(path string, handler Handler)
	POST(path string, handler Handler)
	PUT(path string, handler Handler)
	DELETE(path string, handler Handler)
	Use(middlewares ...Middleware)
	Shutdown(ctx context.Context) error
	Start(addr string) error
	Group(prefix string, middlewares ...Middleware) RouterGroup
}

type RouterGroup interface {
	Origin() any
	GET(path string, handler Handler)
	POST(path string, handler Handler)
	PUT(path string, handler Handler)
	DELETE(path string, handler Handler)
	Use(middlewares ...Middleware)
	Group(prefix string, middlewares ...Middleware) RouterGroup
}

type Response struct {
	Code    string `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func applyMiddleware(h Handler, middlewares []Middleware) Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func copyMiddlewares(middlewares []Middleware, ap ...Middleware) []Middleware {
	return append([]Middleware{}, middlewares...)
}
