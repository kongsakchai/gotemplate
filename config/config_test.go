package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("should load config with default values", func(t *testing.T) {
		Env = "LOCAL" // Reset environment variable for testing

		os.Clearenv()
		t.Setenv("APP_NAME", "TestApp")
		t.Setenv("APP_PORT", "8080")
		t.Setenv("APP_VERSION", "1.0.0")
		t.Setenv("DATABASE_URL", "Use Here")
		t.Setenv("LOCAL_DATABASE_URL", "Not Use Here")
		t.Setenv("LOCAL_REDIS_HOST", "Use Here")

		expectConfig := Config{
			App: App{
				Name:    "TestApp",
				Port:    "8080",
				Version: "1.0.0",
			},
			Header: Header{
				RefIDKey: "X-Ref-ID",
			},
			Database: Database{
				URL: "Use Here",
			},
			Migration: Migration{
				Directory: "./migrations",
			},
			Redis: Redis{
				Host:    "Use Here",
				Timeout: time.Minute * 10,
			},
		}

		cfg := Load(Env)
		assert.Equal(t, expectConfig, cfg)
	})
}
