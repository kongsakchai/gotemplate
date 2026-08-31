package app

import "github.com/labstack/echo/v5"

func ErrorMiddleware(handler func(error) Error) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}
			if appErr, ok := err.(Error); ok {
				return appErr
			}
			return handler(err)
		}
	}
}
