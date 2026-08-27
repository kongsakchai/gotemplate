package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMySQL(t *testing.T) {
	// MySQL
	t.Run("should error when mysql connection fail", func(t *testing.T) {
		// check panic
		defer func() {
			p := recover()
			assert.NotNil(t, p)
		}()

		New("root:example@(localhost:1111)/example")
	})
}
