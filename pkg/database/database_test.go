package database

import (
	"testing"

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

		NewDatabase("unknow", "invalid")
	})

	t.Run("should ping success when db connet", func(t *testing.T) {
		// check panic
		defer func() {
			p := recover()
			assert.Nil(t, p)
		}()

		db, close := NewDatabase("sqlite", ":memory:")
		defer close(t.Context())

		err := db.Ping()

		assert.NoError(t, err)
	})
}
