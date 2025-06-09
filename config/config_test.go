package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	Env = "" // Reset environment variable for testing

	t.Setenv("APP_NAME", "TestApp")
	t.Setenv("APP_PORT", "8080")
	t.Setenv("APP_VERSION", "1.0.0")
	t.Setenv("HEADER_REF_ID_KEY", "X-Ref-ID")
	t.Setenv("DATABASE_URL", "user:password@tcp(localhost:5432)/testdb")
	t.Setenv("REDIS_HOST", "localhost")
	t.Setenv("REDIS_PORT", "6379")
	t.Setenv("REDIS_USERNAME", "redisuser")
	t.Setenv("REDIS_PASSWORD", "redispassword")
	t.Setenv("REDIS_DB", "1")
	t.Setenv("REDIS_TIMEOUT", "5s")
	t.Setenv("MIGRATION_DIRECTORY", "./migrations")
	t.Setenv("MIGRATION_VERSION", "")

	expectConfig := Config{
		App: App{
			Name:    "TestApp",
			Port:    "8080",
			Version: "1.0.0",
		},
		Header: Header{
			RefIDKey: "X-Ref-ID",
		},
		Migration: Migration{
			Directory: "./migrations",
			Version:   "",
		},
		Database: Database{
			URL: "user:password@tcp(localhost:5432)/testdb",
		},
		Redis: Redis{
			Host:     "localhost",
			Port:     "6379",
			Username: "redisuser",
			Password: "redispassword",
			DB:       1,
			Timeout:  5 * time.Second,
		},
	}

	cfg := Load()

	assert.Equal(t, expectConfig, cfg)
}
