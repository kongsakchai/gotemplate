package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func prefix(key string) string {
	if Env == "" {
		return key
	}

	return fmt.Sprintf("%s_%s", Env, key)
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	val := os.Getenv(prefix(key))
	if val == "" {
		return defaultValue
	}

	dur, err := time.ParseDuration(val)
	if err != nil {
		return defaultValue
	}
	return dur
}

func getInt(key string, defaultValue int) int {
	val := os.Getenv(prefix(key))
	if val == "" {
		return defaultValue
	}

	intVal, err := strconv.Atoi(val)
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
		RefIDKey: os.Getenv("HEADER_REF_ID_KEY"),
	}
}

func getDatabaseConfig() Database {
	return Database{
		Host:     os.Getenv(prefix("DB_HOST")),
		Port:     os.Getenv(prefix("DB_PORT")),
		Username: os.Getenv(prefix("DB_USERNAME")),
		Password: os.Getenv(prefix("DB_PASSWORD")),
		DBName:   os.Getenv(prefix("DB_NAME")),
	}
}

func getRedisConfig() Redis {
	return Redis{
		Host:     os.Getenv(prefix("REDIS_HOST")),
		Port:     os.Getenv(prefix("REDIS_PORT")),
		Username: os.Getenv(prefix("REDIS_USERNAME")),
		Password: os.Getenv(prefix("REDIS_PASSWORD")),
		DB:       getInt(prefix("REDIS_DB"), 0),
		Timeout:  getDuration(prefix("REDIS_TIMEOUT"), 5*time.Second),
	}
}
