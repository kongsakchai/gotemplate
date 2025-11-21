package cache

import (
	"github.com/kongsakchai/gotemplate/config"
	redis "github.com/redis/go-redis/v9"
)

func NewRedis(conf config.Redis) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         conf.Host + ":" + conf.Port,
		Username:     conf.Username,
		Password:     conf.Password,
		DB:           conf.DB,
		ReadTimeout:  conf.Timeout,
		WriteTimeout: conf.Timeout,
	})
}
