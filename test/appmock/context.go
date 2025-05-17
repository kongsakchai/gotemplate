package appmock

import (
	"context"
	"log/slog"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/stretchr/testify/mock"
)

type Context struct {
	mock.Mock
	app.Context
}

func (c *Context) Query(key string) string {
	args := c.Called(key)
	return args.String(0)
}

func (c *Context) Param(key string) string {
	args := c.Called(key)
	return args.String(0)
}

func (c *Context) Bind(obj any) error {
	args := c.Called(obj)
	return args.Error(0)
}

func (c *Context) JSON(code int, obj any) error {
	args := c.Called(code, obj)
	return args.Error(0)
}

func (c *Context) OK(obj any) error {
	args := c.Called(obj)
	return args.Error(0)
}

func (c *Context) OKWithMessage(message string, obj any) error {
	args := c.Called(message, obj)
	return args.Error(0)
}

func (c *Context) Created(obj any) error {
	args := c.Called(obj)
	return args.Error(0)
}

func (c *Context) CreatedWithMessage(message string, obj any) error {
	args := c.Called(message, obj)
	return args.Error(0)
}

func (c *Context) NotFound(err *app.Error) error {
	args := c.Called(err)
	return args.Error(0)
}

func (c *Context) InternalServer(err *app.Error) error {
	args := c.Called(err)
	return args.Error(0)
}

func (c *Context) BadRequest(err *app.Error) error {
	args := c.Called(err)
	return args.Error(0)
}

func (c *Context) Ctx() context.Context {
	args := c.Called()
	return args.Get(0).(context.Context)
}

func (c *Context) Get(key string) any {
	args := c.Called(key)
	return args.Get(0)
}

func (c *Context) Set(key string, value any) {
	c.Called(key, value)
}

func (c *Context) Logger() *slog.Logger {
	args := c.Called()
	return args.Get(0).(*slog.Logger)
}
