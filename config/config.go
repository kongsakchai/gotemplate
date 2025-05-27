package config

import (
	"sync"
)

type Config struct {
	App    App
	Header Header
}

type App struct {
	Name    string
	Port    string
	Version string
}

type Header struct {
	RefIDKey string
}

var config Config
var once sync.Once

func Load() Config {
	once.Do(func() {
		config = Config{
			App:    getAppConfig(),
			Header: getHeaderConfig(),
		}
	})

	return config
}
