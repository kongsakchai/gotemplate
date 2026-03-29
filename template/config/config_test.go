package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("should load config with default values", func(t *testing.T) {
		Env = "LOCAL" // Reset environment variable for testing

		os.Clearenv()
		t.Setenv("APP_NAME", "TestApp")
		t.Setenv("APP_PORT", "8080")
		t.Setenv("APP_VERSION", "1.0.0")

		expectConfig := App{
			Name:    "TestApp",
			Port:    "8080",
			Version: "1.0.0",
		}

		cfg := Load(Env)
		assert.Equal(t, expectConfig, cfg.App)
	})
}
