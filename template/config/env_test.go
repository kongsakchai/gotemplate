package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	t.Run("should return true for Local environment", func(t *testing.T) {
		Env = Local
		assert.True(t, IsLocal())
	})

	t.Run("should return true for Dev environment", func(t *testing.T) {
		Env = Dev
		assert.True(t, IsDev())
	})

	t.Run("should return true for Prod environment", func(t *testing.T) {
		Env = Prod
		assert.True(t, IsProd())
	})
}
