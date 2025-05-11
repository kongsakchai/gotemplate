package config

import (
	"os"
	"strconv"
	"sync"
	"time"
)

var once sync.Once
var conf Config

func Load() Config {
	once.Do(func() {
		conf = Config{
			App:      getApp(),
			Database: getDatabase(),
			Redis:    getRedis(),
		}
	})

	return conf
}

func getDatabase() Database {
	return Database{
		MySQLURI: getStr("DATABASE_MYSQL_URI", ""),
	}
}

func getRedis() Redis {
	return Redis{
		Host:           getStr("REDIS_HOST", ""),
		Port:           getStr("REDIS_PORT", ""),
		Password:       getStr("REDIS_PASSWORD", ""),
		DB:             getInt("REDIS_DB", 0),
		Timeout:        getDuration("REDIS_TIMEOUT", 5*time.Second),
		ConnectTimeout: getDuration("REDIS_CONNECT_TIMEOUT", 5*time.Second),
	}
}

func getApp() App {
	return App{
		Service:  getStr("APP_SERVICE", "myapp"),
		Port:     getStr("APP_PORT", "8080"),
		Env:      getStr("APP_ENV", "DEV"),
		LogLevel: getStr("APP_LOG_LEVEL", "info"),
	}
}

func getStr(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return duration
}
