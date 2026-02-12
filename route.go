package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/config"
	"github.com/kongsakchai/gotemplate/database"
	"github.com/kongsakchai/gotemplate/middleware"
	"github.com/kongsakchai/gotemplate/validator"
	"github.com/labstack/echo/v4"
)

func router(cfg config.Config) (app.App, []shutdownFunc) {
	db, closeDB := database.NewMySQL(cfg.Database)
	setMigration(db.DB, cfg.Migration)

	r := app.NewEchoApp()
	r.Validator = validator.NewReqValidator()

	r.Use(
		middleware.RefID(cfg.Header.RefIDKey),
		middleware.Logger(cfg.Header.RefIDKey, true, true),
	)

	r.GET("/health", healthCheck(db))

	return r, []shutdownFunc{
		r.Shutdown,
		closeDB,
	}
}

func healthCheck(db *sqlx.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if db.Ping() != nil {
			return app.Fail(ctx, app.InternalServer(app.ErrCode, app.ErrDatabaseMsg, nil))
		}
		return app.OkWithMessage(ctx, nil, "healthy")
	}
}
