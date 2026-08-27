package postgresdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPostgres(t *testing.T) {
	// Postgres
	t.Run("should error when postgres connection fail", func(t *testing.T) {
		// check panic
		defer func() {
			p := recover()
			assert.NotNil(t, p)
		}()

		New("root:example@(localhost:1111)/example")
	})
}
