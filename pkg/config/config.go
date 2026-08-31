package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	App       App
	Header    Header
	Migration Migration
	Database  Database
	Redis     Redis
	Log       Log
	AppJWT    JWT
}

type App struct {
	Name    string `env:"APP_NAME" envDefault:"gotemplate"`
	Port    string `env:"APP_PORT" envDefault:"8080"`
	Version string `env:"APP_VERSION" envDefault:"0.0.1"`
}

type Header struct {
	RefIDKey string `env:"HEADER_REF_ID_KEY" envDefault:"X-Ref-ID"`
}

type Migration struct {
	Enable      bool   `env:"MIGRATION_ENABLE" envDefault:"false"`
	Source      string `env:"MIGRATION_SRC" envDefault:"file://migrations"`
	DatabaseURL string `env:"MIGRATION_DATABASE_URL"`
	Version     int    `env:"MIGRATION_VERSION" envDefault:"-1"`
}

type Database struct {
	URL string `env:"DATABASE_URL"`
}

type Redis struct {
	Host     string        `env:"REDIS_HOST"`
	Port     string        `env:"REDIS_PORT"`
	Username string        `env:"REDIS_USERNAME"`
	Password string        `env:"REDIS_PASSWORD"`
	DB       int           `env:"REDIS_DB" envDefault:"0"`
	Timeout  time.Duration `env:"REDIS_TIMEOUT" envDefault:"10m"`
}

type Log struct {
	Enable     bool `env:"LOG_ENABLE"`
	HttpEnable bool `env:"LOG_HTTP_ENABLE"`
}

type JWT struct {
	Method    string        `env:"JWT_METHOD"`
	SecretKey string        `env:"JWT_SECRET_KEY"`
	VerifyKey string        `env:"JWT_VERIFY_KEY"`
	Issuer    string        `env:"JWT_ISSUER"`
	Audience  string        `env:"JWT_AUDIENCE"`
	Expired   time.Duration `env:"JWT_EXPIRED" envDefault:"15m"`
}

var config Config
var once sync.Once

func Load(prefix string) Config {
	once.Do(func() {
		env.Parse(&config)
		opt := env.Options{
			Prefix:                       fmt.Sprintf("%s_", prefix),
			SetDefaultsForZeroValuesOnly: true,
		}
		env.ParseWithOptions(&config, opt)
	})

	return config
}
