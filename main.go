package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/app/authapp"
	"github.com/kongsakchai/gotemplate/app/todoapp"
	"github.com/kongsakchai/gotemplate/pkg/config"
	"github.com/kongsakchai/gotemplate/pkg/database/sqlitedb"
	"github.com/kongsakchai/gotemplate/pkg/hash"
	"github.com/kongsakchai/gotemplate/pkg/jwttoken"
	"github.com/kongsakchai/gotemplate/pkg/logger"
	"github.com/kongsakchai/gotemplate/pkg/migrate"
	"github.com/labstack/echo/v5"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
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

	// init dependencies that need to be closed after use
	db, close := sqlitedb.New(cfg.Database.URL)
	defer close(context.Background())
	migrate.Migrate(cfg.Migration)

	echoApp := app.NewEchoApp(cfg)
	echoApp.Logger = logger
	registerRoutes(cfg, echoApp, db)

	runApp(echoApp, cfg, gracefulTimeout)
}

func registerRoutes(cfg config.Config, echo *app.EchoApp, db *sqlx.DB) {
	// init dependencies that are not closed after use
	jwt := jwttoken.NewJWTToken(cfg.AppJWT)
	hasher := hash.NewHasher(12)

	app.GET(echo, "health", "/health", healthCheck(db))
	app.GET(echo, "metrics", "/metrics", metrics())

	authapp.NewApp(authapp.Deps{DB: db, Hasher: hasher, Signer: jwt}).RegisterRoute(echo)
	todoapp.NewApp(todoapp.Deps{DB: db, Verifier: jwt}).RegisterRoutes(echo)
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
