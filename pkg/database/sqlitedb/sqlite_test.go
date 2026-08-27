package sqlitedb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSQLite(t *testing.T) {
	// Postgres
	t.Run("should error when sqlite connection fail", func(t *testing.T) {
		// check panic
		defer func() {
			p := recover()
			assert.NotNil(t, p)
		}()

		New("root:example@(localhost:1111)/example")
	})

	t.Run("should success when sqlite with memory success", func(t *testing.T) {
		db, close := NewWithMemory()
		defer close(t.Context())

		err := db.Ping()

		assert.NoError(t, err)
	})
}
