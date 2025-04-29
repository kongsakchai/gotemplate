package config

import (
	"os"
	"sync"
)

var once sync.Once
var conf Config

func Load() Config {
	once.Do(func() {
		conf = Config{
			App:      getApp(),
			Database: getDatabase(),
		}
	})

	return conf
}

func getDatabase() Database {
	return Database{
		MySQLURI: getEnv("DATABASE_MYSQL_URI", ""),
	}
}

func getApp() App {
	return App{
		Service:  getEnv("APP_SERVICE", "myapp"),
		Port:     getEnv("APP_PORT", "8080"),
		Env:      getEnv("APP_ENV", "DEV"),
		LogLevel: getEnv("APP_LOG_LEVEL", "info"),
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
