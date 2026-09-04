package app

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type Router interface {
	AddRoute(routable echo.Route) (echo.RouteInfo, error)
}

func POST(router Router, name string, path string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) {
	router.AddRoute(echo.Route{
		Method:      http.MethodPost,
		Path:        path,
		Handler:     handler,
		Middlewares: middlewares,
		Name:        name,
	})
}

func GET(router Router, name string, path string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) {
	router.AddRoute(echo.Route{
		Method:      http.MethodGet,
		Path:        path,
		Handler:     handler,
		Middlewares: middlewares,
		Name:        name,
	})
}

func PUT(router Router, name string, path string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) {
	router.AddRoute(echo.Route{
		Method:      http.MethodPut,
		Path:        path,
		Handler:     handler,
		Middlewares: middlewares,
		Name:        name,
	})
}

func DELETE(router Router, name string, path string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) {
	router.AddRoute(echo.Route{
		Method:      http.MethodDelete,
		Path:        path,
		Handler:     handler,
		Middlewares: middlewares,
		Name:        name,
	})
}

func PATCH(router Router, name string, path string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) {
	router.AddRoute(echo.Route{
		Method:      http.MethodPatch,
		Path:        path,
		Handler:     handler,
		Middlewares: middlewares,
		Name:        name,
	})
}
