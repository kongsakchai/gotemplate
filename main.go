package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/template/app"
	"github.com/kongsakchai/gotemplate/template/app/member"
	"github.com/kongsakchai/gotemplate/template/pkg/clock"
	"github.com/kongsakchai/gotemplate/template/pkg/config"
	"github.com/kongsakchai/gotemplate/template/pkg/database"
	"github.com/kongsakchai/gotemplate/template/pkg/logger"
	migrate "github.com/kongsakchai/simple-sql-migrate"
	"github.com/labstack/echo/v5"
)

const gracefulTimeout = time.Second * 10

const (
	Byte uint64 = 1 << (10 * iota)
	KB
	MB
)

func init() {
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(1) // 0 - 999m
	}
	if os.Getenv("GOMEMLIMIT") != "" {
		debug.SetMemoryLimit(-1) // GOMEMLIMIT
	}
}

func main() {
	logger := logger.New()
	cfg := config.Load(config.Env)

	db, close := database.NewMySQL(cfg.Database.URL)
	defer close(context.Background())

	clock := clock.New()

	app := app.NewEchoApp(cfg)
	app.Logger = logger

	app.GET("/health", healthCheck(nil))
	app.GET("/metrics", metrics())

	memberMo := member.NewModule(member.External{DB: db, Clock: clock})
	memberMo.Handler.RegisterMemberHandler(app)

	runApp(app, cfg, gracefulTimeout)
}

func runApp(app *app.EchoApp, cfg config.Config, gracefulTimeout time.Duration) {
	slog.Info(cfg.App.Name, "version", cfg.App.Version, "env", config.Env)
	slog.Info("listening on port " + cfg.App.Port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	if err := app.Start(ctx, fmt.Sprintf(":%s", cfg.App.Port), gracefulTimeout); err != nil && err != http.ErrServerClosed {
		slog.Error("shutting down the server: " + err.Error())
		return
	}

	slog.Info("bye bye")
}

func migrateDB(db *sql.DB, cfg config.Migration) {
	if !cfg.Enable {
		return
	}
	if err := migrate.Migrate(db, migrate.Options{
		Source:  cfg.Directory,
		Version: cfg.Version,
		Repeat:  migrate.GetRepeatAction(cfg.Repeat),
	}); err != nil {
		panic("migration failed: " + err.Error())
	}

	version, _ := migrate.Version(db)
	slog.Info("migration completed", "version", version)
}

func healthCheck(db *sqlx.DB) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		if db != nil && db.Ping() != nil {
			return app.Fail(ctx, app.InternalError(app.DatabaseNotReadyCode, app.DatabaseNotReadyMsg, nil))
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
