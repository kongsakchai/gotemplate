package main

import (
	"fmt"
	"runtime"

	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/app/apperror"
	"github.com/kongsakchai/gotemplate/app/middleware"
	"github.com/kongsakchai/gotemplate/config"
	"github.com/kongsakchai/gotemplate/database"
	"github.com/kongsakchai/gotemplate/validator"
	"github.com/labstack/echo/v4"
)

const (
	Byte uint64 = 1 << (10 * iota)
	KB
	MB
)

func router(cfg config.Config) (app.App, []shutdownFunc) {
	db, closeDB := database.NewMySQL(cfg.Database)
	setMigration(db.DB, cfg.Migration)

	r := app.NewEchoApp()
	r.Validator = validator.NewReqValidator()
	r.HTTPErrorHandler = apperror.ErrorHandler

	r.Use(
		middleware.RefID(cfg.Header.RefIDKey),
		middleware.Logger(true),
	)

	r.GET("/health", healthCheck(db))
	r.GET("/metrics", metrics())

	return r, []shutdownFunc{
		r.Shutdown,
		closeDB,
	}
}

func healthCheck(db *sqlx.DB) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if db.Ping() != nil {
			return app.Fail(ctx, app.InternalError(app.ErrInternalCode, app.ErrDatabaseMsg, nil))
		}
		return app.Ok(ctx, nil, "healthy")
	}
}

func toMB(b uint64) string {
	return fmt.Sprintf("%.2f MB", float64(b)/float64(MB))
}

func metrics() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		return app.Ok(ctx, map[string]string{
			"alloc":        toMB(mem.Alloc),
			"totalAlloc":   toMB(mem.TotalAlloc),
			"sysAlloc":     toMB(mem.Sys),
			"heapInuse":    toMB(mem.HeapInuse),
			"heapIdle":     toMB(mem.HeapIdle),
			"heapReleased": toMB(mem.HeapReleased),
			"stackInuse":   toMB(mem.StackInuse),
			"stackSys":     toMB(mem.StackSys),
		})
	}
}
