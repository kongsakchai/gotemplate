package cache

import (
	"time"

	redis "github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DB       int
	Timeout  time.Duration
}

func NewRedis(conf RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         conf.Host + ":" + conf.Port,
		Username:     conf.Username,
		Password:     conf.Password,
		DB:           conf.DB,
		ReadTimeout:  conf.Timeout,
		WriteTimeout: conf.Timeout,
	})
}
