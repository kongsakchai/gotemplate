package config

import (
	"fmt"
	"os"
)

func prefix(key string) string {
	if Env == "" {
		return key
	}

	return fmt.Sprintf("%s_%s", Env, key)
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
