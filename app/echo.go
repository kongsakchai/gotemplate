package app

import (
	"context"
	"log/slog"

	"github.com/kongsakchai/gotemplate/pkg/generate"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	echoContextPool *pool[EchoContext]
)

func init() {
	echoContextPool = createPool[EchoContext]()
}

type EchoContext struct {
	logger *slog.Logger
	echo.Context
}

func (e *EchoContext) reset(ctx echo.Context, logger *slog.Logger) {
	e.Context = ctx
	e.logger = logger
}

func (e *EchoContext) Query(key string) string {
	return e.Context.QueryParam(key)
}

func (e *EchoContext) Param(key string) string {
	return e.Context.Param(key)
}

func (e *EchoContext) Bind(obj any) error {
	return e.Context.Bind(obj)
}

func (e *EchoContext) JSON(code int, obj any) error {
	return e.Context.JSON(code, obj)
}

func (e *EchoContext) Validate(obj any) error {
	return e.Context.Validate(obj)
}

func (e *EchoContext) OK(obj any) error {
	return e.JSON(200, Response{
		Status: SuccessStatus,
		Code:   SuccessCode,
		Data:   obj,
	})
}

func (e *EchoContext) OKWithMessage(message string, obj any) error {
	return e.JSON(200, Response{
		Status:  SuccessStatus,
		Code:    SuccessCode,
		Message: message,
		Data:    obj,
	})
}

func (e *EchoContext) Created(obj any) error {
	return e.JSON(201, Response{
		Status: SuccessStatus,
		Code:   SuccessCode,
		Data:   obj,
	})
}

func (e *EchoContext) CreatedWithMessage(message string, obj any) error {
	return e.JSON(201, Response{
		Status:  SuccessStatus,
		Code:    SuccessCode,
		Message: message,
		Data:    obj,
	})
}

func (e *EchoContext) Error(err Error) error {
	e.logger.Error(err.Error())
	return e.JSON(err.StatusCd, Response{
		Status:  ErrorStatus,
		Code:    err.Code,
		Message: err.Message,
	})
}

func (e *EchoContext) Ctx() context.Context {
	return e.Context.Request().Context()
}

func (e *EchoContext) Get(key string) any {
	return e.Context.Get(key)
}

func (e *EchoContext) Set(key string, value any) {
	e.Context.Set(key, value)
}

func (e *EchoContext) Logger() *slog.Logger {
	return e.logger
}

type echoRouter struct {
	*echo.Echo
	logger *slog.Logger
}

func NewEchoRoute(logger *slog.Logger) *echoRouter {
	e := echo.New()
	e.HideBanner = true

	// e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
	}))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("traceID", generate.UUID())
			return next(c)
		}
	})

	return &echoRouter{
		Echo:   e,
		logger: logger,
	}
}

func (e *echoRouter) NewHandler(handler Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := echoContextPool.Get()
		defer echoContextPool.Put(ctx)
		ctx.reset(c, e.logger)
		return handler(ctx)
	}
}
