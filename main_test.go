package main

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kongsakchai/gotemplate/config"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

func TestSetupRoutes(t *testing.T) {
	t.Run("should return healthy when can ping db success", func(t *testing.T) {
		t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
		// init container
		ct, err := testcontainers.Run(
			t.Context(),
			"mariadb:latest",

			testcontainers.WithProvider(testcontainers.ProviderPodman),
			testcontainers.WithExposedPorts("3306/tcp"),
			testcontainers.WithWaitStrategy(
				wait.ForListeningPort("3306/tcp"),
			),

			testcontainers.WithEnv(map[string]string{
				"MYSQL_ROOT_PASSWORD": "example",
				"MYSQL_DATABASE":      "example",
			}),
		)
		require.NoError(t, err)

		// clear container
		defer testcontainers.CleanupContainer(t, ct)

		endpoint, err := ct.Endpoint(t.Context(), "")
		require.NoError(t, err)

		// arrange
		app, shutdown := router(config.Config{
			Database: config.Database{
				URL: fmt.Sprintf("root:example@(%s)/example", endpoint),
			},
		}, slog.Default())
		defer shutdown(context.Background())

		go func() {
			err := app.Start(t.Context(), ":8888")
			assert.Error(t, http.ErrServerClosed, err)
		}()
		time.Sleep(1 * time.Second)

		// act
		resp, err := http.Get("http://localhost:8888/health")
		assert.NoError(t, err)

		//assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		app.Shutdown(t.Context())
	})
}

func TestHealthCheck(t *testing.T) {
	t.Run("should return healthy when can ping db success", func(t *testing.T) {
		db, err := sqlx.Open("sqlite", ":memory:")
		require.NoError(t, err)
		defer db.Close()

		// arrange
		ctx, rec := echotest.ContextConfig{}.ToContextRecorder(t)

		handler := healthCheck(db)

		// act
		handler(ctx)

		//assert
		assert.NotNil(t, rec.Body)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should return error when can ping db fail", func(t *testing.T) {
		db, err := sqlx.Open("mysql", "root:example@(localhost:0000)/example")
		require.NoError(t, err)
		defer db.Close()

		// arrange
		ctx, rec := echotest.ContextConfig{}.ToContextRecorder(t)

		handler := healthCheck(db)

		// act
		handler(ctx)

		//assert
		assert.NotNil(t, rec.Body)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

type mockApp struct {
	startErr    error
	shutdownErr error
}

func (m *mockApp) Start(ctx context.Context, addr string) error {
	<-ctx.Done()
	return m.startErr
}

func (m *mockApp) Shutdown(ctx context.Context) error {
	<-ctx.Done()
	return m.shutdownErr
}

func TestGracefulShutdown(t *testing.T) {
	d := gracefulTimeout
	gracefulTimeout = 1 * time.Second
	t.Run("should not panic when graceful shutdown success", func(t *testing.T) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
		defer stop()

		idle := make(chan struct{})
		go func() {
			// check panic
			defer func() {
				p := recover()
				assert.Nil(t, p)
			}()

			gracefulShutdown(idle, ctx, &mockApp{})
		}()

		time.Sleep(1 * time.Second)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		start := time.Now()
		<-idle
		dur := time.Since(start)

		assert.Equal(t, gracefulTimeout.Seconds(), math.Round(dur.Seconds()))
	})

	t.Run("should panic when graceful shutdown failed", func(t *testing.T) {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
		defer stop()
		idle := make(chan struct{})
		go func() {
			// check panic
			defer func() {
				p := recover()
				assert.NotNil(t, p)
			}()

			gracefulShutdown(idle, ctx, &mockApp{shutdownErr: fmt.Errorf("failed to shutdown")})
		}()

		time.Sleep(1 * time.Second)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		<-idle
	})

	gracefulTimeout = d
}

func TestSetMigration(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	{
		// act
		setMigration(db.DB, config.Migration{Enable: true, Version: "0000", Directory: "./migrations/mock"})
		// assert
		data, err := db.QueryContext(t.Context(), "SELECT * FROM mock_data")
		assert.NoError(t, err)
		assert.False(t, data.Next())
	}

	{
		// act
		setMigration(db.DB, config.Migration{Enable: true, Version: "", Directory: "./migrations/mock"})
		// assert
		data, err := db.QueryContext(t.Context(), "SELECT * FROM mock_data")
		assert.NoError(t, err)
		assert.True(t, data.Next())

		var id int
		var val int
		data.Scan(&id, &val)

		assert.Equal(t, 1, id)
		assert.Equal(t, 10, val)
	}

	{
		setMigration(db.DB, config.Migration{Enable: false, Directory: "invalid"})
	}

	{
		defer func() {
			p := recover()
			assert.NotNil(t, p)
		}()
		setMigration(db.DB, config.Migration{Enable: true, Directory: "invalid"})
	}

}
