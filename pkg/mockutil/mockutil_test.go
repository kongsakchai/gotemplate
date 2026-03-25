package mockutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimer(t *testing.T) {
	t.Run("should create a new timer", func(t *testing.T) {
		// act
		timer := NewTimer(t)

		// assert
		assert.NotNil(t, timer)
	})
}
