package database

import (
	"fmt"
	"testing"

	"github.com/kongsakchai/gotemplate/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestNewDatabase(t *testing.T) {
	t.Run("should panic when unknow driver", func(t *testing.T) {
		// check panic
		defer func() {
			p := recover()
			assert.NotNil(t, p)
		}()

		newDatabase("unknow", "invalid")
	})

	// MySQL
	t.Run("should ping panic when not found db", func(t *testing.T) {
		// check panic
		defer func() {
			p := recover()
			assert.NotNil(t, p)
		}()

		NewMySQL(config.Database{
			URL: "root:example@(localhost:mock)/example2",
		})
	})

	t.Run("should ping no error when mysql connection success", func(t *testing.T) {
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

		// check panic
		defer func() {
			p := recover()
			assert.Nil(t, p)
		}()

		db, close := NewMySQL(config.Database{
			URL: fmt.Sprintf("root:example@(%s)/example", endpoint),
		})

		assert.NoError(t, db.Ping())
		close()
		assert.Error(t, db.Ping())

		testcontainers.CleanupContainer(t, ct)
	})

	// Postgres
	t.Run("should ping no error when connection success", func(t *testing.T) {
		t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
		// init container
		ct, err := testcontainers.Run(
			t.Context(),
			"postgres:latest",

			testcontainers.WithProvider(testcontainers.ProviderPodman),
			testcontainers.WithExposedPorts("5432/tcp"),
			testcontainers.WithWaitStrategy(
				wait.ForListeningPort("5432/tcp"),
			),

			testcontainers.WithEnv(map[string]string{
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "example",
				"POSTGRES_DB":       "example",
			}),
		)
		require.NoError(t, err)

		endpoint, err := ct.Endpoint(t.Context(), "")
		require.NoError(t, err)

		// check panic
		defer func() {
			p := recover()
			assert.Nil(t, p)
		}()

		db, close := NewPostgres(config.Database{
			URL: fmt.Sprintf("postgres://postgres:example@%s/example?sslmode=disable", endpoint),
		})

		assert.NoError(t, db.Ping())
		close()
		assert.Error(t, db.Ping())

		// clear container
		testcontainers.CleanupContainer(t, ct)
	})
}
