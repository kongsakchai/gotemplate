package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func prefix(key string) string {
	if Env == "" {
		return key
	}

	return fmt.Sprintf("%s_%s", Env, key)
}

func getSecret(key string) string {
	if !strings.HasPrefix(key, "$") {
		return key
	}
	secret := os.ExpandEnv(key)
	if secret == "" {
		return key
	}
	return secret
}

func getString(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	return getSecret(val)
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	dur, err := time.ParseDuration(getSecret(val))
	if err != nil {
		return defaultValue
	}
	return dur
}

func getInt(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	intVal, err := strconv.Atoi(getSecret(val))
	if err != nil {
		return defaultValue
	}
	return intVal
}

func getAppConfig() App {
	return App{
		Name:    os.Getenv("APP_NAME"),
		Port:    os.Getenv("APP_PORT"),
		Version: os.Getenv("APP_VERSION"),
	}
}

func getHeaderConfig() Header {
	return Header{
		RefIDKey: getString("HEADER_REF_ID_KEY", "X-Ref-ID"),
	}
}

func getMigrationConfig() Migration {
	return Migration{
		Directory: getString("MIGRATION_DIRECTORY", "./migrations"),
		Version:   getString("MIGRATION_VERSION", ""),
		TableName: getString("MIGRATION_TABLE_NAME", "schema_migrations"),
	}
}

func getDatabaseConfig() Database {
	return Database{
		URL: getString(prefix("DATABASE_URL"), ""),
	}
}

func getRedisConfig() Redis {
	return Redis{
		Host:     getString(prefix("REDIS_HOST"), ""),
		Port:     getString(prefix("REDIS_PORT"), ""),
		Username: getString(prefix("REDIS_USERNAME"), ""),
		Password: getString(prefix("REDIS_PASSWORD"), ""),
		DB:       getInt(prefix("REDIS_DB"), 0),
		Timeout:  getDuration(prefix("REDIS_TIMEOUT"), 5*time.Second),
	}
}
