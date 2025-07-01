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
	Directory string `env:"MIGRATION_DIR" envDefault:"./migrations"`
	Version   string `env:"MIGRATION_VERSION"`
	TableName string `env:"MIGRATION_TABLE_NAME"`
}

type Database struct {
	URL string `envKey:"DATABASE_URL"`
}

type Redis struct {
	Host     string        `envKey:"REDIS_HOST"`
	Port     string        `envKey:"REDIS_PORT"`
	Username string        `envKey:"REDIS_USERNAME"`
	Password string        `envKey:"REDIS_PASSWORD"`
	DB       int           `envKey:"REDIS_DB" envKeyDefault:"0"`
	Timeout  time.Duration `envKey:"REDIS_TIMEOUT" envKeyDefault:"10m"`
}

var config Config
var once sync.Once

// TODO: when add new config \n
//   - tag `env` is used for environment variables. envDefault is used for default value.
//   - ex. `env:"DATABASE_URL"` -> `DATABASE_URL`
//   - tag `envKey` is used for environment variables with prefix. envKeyDefault is used for default value.
//   - ex. `envKey:"DATABASE_URL"` -> `XXX_DATABASE_URL`
func Load(prefix string) Config {
	once.Do(func() {
		env.Parse(&config)
		opt := env.Options{
			Prefix:              fmt.Sprintf("%s_", prefix),
			TagName:             "envKey",
			DefaultValueTagName: "envKeyDefault",
		}
		env.ParseWithOptions(&config, opt)
	})

	return config
}
