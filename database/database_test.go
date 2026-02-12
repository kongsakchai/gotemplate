package database

import (
	"testing"

	"github.com/kongsakchai/gotemplate/config"
	"github.com/stretchr/testify/assert"

	_ "modernc.org/sqlite"
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

	t.Run("should ping success when db connet", func(t *testing.T) {
		// check panic
		defer func() {
			p := recover()
			assert.Nil(t, p)
		}()

		db, close := newDatabase("sqlite", ":memory:")
		defer close(t.Context())

		err := db.Ping()

		assert.NoError(t, err)
	})

	// MySQL
	t.Run("should error when mysql connection fail", func(t *testing.T) {
		// check panic
		defer func() {
			p := recover()
			assert.NotNil(t, p)
		}()

		NewMySQL(config.Database{
			URL: "root:example@(localhost:1111)/example",
		})
	})

	// Postgres
	t.Run("should error when mysql connection fail", func(t *testing.T) {
		// check panic
		defer func() {
			p := recover()
			assert.NotNil(t, p)
		}()

		NewPostgres(config.Database{
			URL: "root:example@(localhost:1111)/example",
		})
	})
}
