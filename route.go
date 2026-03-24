package main

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"

	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/app/apperror"
	"github.com/kongsakchai/gotemplate/app/middleware"
	"github.com/kongsakchai/gotemplate/config"
	"github.com/kongsakchai/gotemplate/validator"
	"github.com/labstack/echo/v5"
	"github.com/redis/go-redis/v9"
)

const (
	Byte uint64 = 1 << (10 * iota)
	KB
	MB
)

type externalService struct {
	DB    *sqlx.DB
	Redis *redis.Client
}

func setupExternalService(cfg config.Config) (externalService, func(context.Context) error) {
	// db, closeDB := database.NewM(cfg.Database)
	// rd := cache.NewRedis(cfg.Redis)

	return externalService{}, func(ctx context.Context) error {
		return nil
	}
}

func router(cfg config.Config, external externalService, logger *slog.Logger) app.App {
	r := app.NewEchoApp()
	r.Validator = validator.NewReqValidator()
	r.HTTPErrorHandler = apperror.ErrorHandler
	r.Logger = logger

	r.Use(
		middleware.RefID(cfg.Header.RefIDKey, cfg.Log.Tags),
		middleware.Logger(cfg.Log.Enable),
	)

	r.GET("/health", healthCheck(external.DB))
	r.GET("/metrics", metrics())

	return r
}

func healthCheck(db *sqlx.DB) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		if db != nil && db.Ping() != nil {
			return app.Fail(ctx, app.InternalError(app.ErrInternalCode, app.ErrDatabaseMsg, nil))
		}
		return app.Ok(ctx, nil, "healthy")
	}
}

func toMB(b uint64) string {
	return fmt.Sprintf("%.2f MB", float64(b)/float64(MB))
}

func metrics() echo.HandlerFunc {
	return func(ctx *echo.Context) error {
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
