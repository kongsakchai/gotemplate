package main

import (
	"database/sql"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/config"
	"github.com/labstack/echo/v4"
)

func router(cfg config.Config, r *app.EchoApp, db *sql.DB) *app.EchoApp {
	r.GET("/health", healthCheck(db))

	return r
}

func healthCheck(db *sql.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if db.Ping() != nil {
			return app.Fail(ctx, app.InternalServer(app.ErrCode, app.ErrDatabaseMsg, nil))
		}
		return app.OkWithMessage(ctx, nil, "healthy")
	}
}
