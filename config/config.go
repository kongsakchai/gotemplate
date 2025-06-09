package config

import (
	"sync"
	"time"
)

type Config struct {
	App       App
	Header    Header
	Migration Migration
	Database  Database
	Redis     Redis
}

type App struct {
	Name    string
	Port    string
	Version string
}

type Header struct {
	RefIDKey string
}

type Migration struct {
	Directory string
	Version   string
}

type Database struct {
	URL string
}

type Redis struct {
	Host     string
	Port     string
	Username string
	Password string
	DB       int
	Timeout  time.Duration
}

var config Config
var once sync.Once

func Load() Config {
	once.Do(func() {
		config = Config{
			App:       getAppConfig(),
			Header:    getHeaderConfig(),
			Migration: getMigrationConfig(),
			Database:  getDatabaseConfig(),
			Redis:     getRedisConfig(),
		}
	})

	return config
}
