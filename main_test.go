package main

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
	"time"

	"github.com/kongsakchai/gotemplate/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/go-sql-driver/mysql"
)

func shutdownsAll(ctx context.Context, shutdowns []shutdownFunc) {
	for _, fn := range shutdowns {
		fn(ctx)
	}
}

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

		endpoint, err := ct.Endpoint(t.Context(), "")
		require.NoError(t, err)

		// arrange
		app, shutdown := router(config.Config{
			Database: config.Database{
				URL: fmt.Sprintf("root:example@(%s)/example", endpoint),
			},
		})
		defer shutdownsAll(t.Context(), shutdown)

		go app.Start(":8888")
		time.Sleep(1 * time.Second)

		// act
		resp, err := http.Get("http://localhost:8888/health")
		assert.NoError(t, err)

		//assert
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		app.Shutdown(t.Context())

		// clear container
		testcontainers.CleanupContainer(t, ct)
	})
}

func TestHealthCheck(t *testing.T) {
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

		endpoint, err := ct.Endpoint(t.Context(), "")
		require.NoError(t, err)

		db, err := sql.Open("mysql", fmt.Sprintf("root:example@(%s)/example", endpoint))
		require.NoError(t, err)

		// arrange
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		handler := healthCheck(db)

		// act
		handler(ctx)

		//assert
		assert.NotNil(t, rec.Body)
		assert.Equal(t, http.StatusOK, rec.Code)

		// clear container
		testcontainers.CleanupContainer(t, ct)
	})

	t.Run("should return error when can ping db fail", func(t *testing.T) {
		db, err := sql.Open("mysql", "root:example@(localhost:3306)/example")
		require.NoError(t, err)

		// arrange
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
		rec := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, rec)

		handler := healthCheck(db)

		// act
		handler(ctx)

		//assert
		assert.NotNil(t, rec.Body)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGracefulShutdown(t *testing.T) {
	d := gracefulTimeout
	gracefulTimeout = 1 * time.Second
	t.Run("should not panic when graceful shutdown success", func(t *testing.T) {
		idle := make(chan struct{})
		go func() {
			// check panic
			defer func() {
				p := recover()
				assert.Nil(t, p)
			}()

			gracefulShutdown(idle, func(ctx context.Context) error {
				<-ctx.Done()
				return nil
			})
		}()

		time.Sleep(1 * time.Second)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		start := time.Now()
		<-idle
		dur := time.Since(start)

		assert.Equal(t, gracefulTimeout.Seconds(), math.Round(dur.Seconds()))
	})

	t.Run("should panic when graceful shutdown failed", func(t *testing.T) {
		idle := make(chan struct{})
		go func() {
			// check panic
			defer func() {
				p := recover()
				assert.NotNil(t, p)
			}()

			gracefulShutdown(idle, func(ctx context.Context) error {
				return fmt.Errorf("force shutdown")
			})
		}()

		time.Sleep(1 * time.Second)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		<-idle
	})

	gracefulTimeout = d
}

func TestSetMigration(t *testing.T) {
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

	endpoint, err := ct.Endpoint(t.Context(), "")
	require.NoError(t, err)

	db, err := sql.Open("mysql", fmt.Sprintf("root:example@(%s)/example", endpoint))
	require.NoError(t, err)

	{
		// act
		setMigration(db, config.Migration{Enable: true, Version: "0000", Directory: "./migrations/mock"})
		// assert
		data, err := db.QueryContext(t.Context(), "SELECT * FROM mock_data")
		assert.NoError(t, err)
		assert.False(t, data.Next())
	}

	{
		// act
		setMigration(db, config.Migration{Enable: true, Version: "", Directory: "./migrations/mock"})
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
		setMigration(db, config.Migration{Enable: false, Directory: "invalid"})
	}

	{
		defer func() {
			p := recover()
			assert.NotNil(t, p)
		}()
		setMigration(db, config.Migration{Enable: true, Directory: "invalid"})
	}

	testcontainers.CleanupContainer(t, ct)
}
