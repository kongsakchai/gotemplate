package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kongsakchai/gotemplate/config"
	redis "github.com/redis/go-redis/v9"
)

type Cache struct {
	client  *redis.Client
	timeout time.Duration
}

func NewCache(conf *config.Redis) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Host + ":" + conf.Port,
		Username: conf.Username,
		Password: conf.Password,
		DB:       conf.DB,
	})

	return &Cache{
		client:  client,
		timeout: conf.Timeout,
	}
}

func (c *Cache) Set(key string, value any, expire time.Duration) error {
	ctx, cancle := context.WithTimeout(context.Background(), c.timeout)
	defer cancle()

	return c.client.Set(ctx, key, value, expire).Err()
}

func (c *Cache) GetStr(key string) (string, error) {
	ctx, cancle := context.WithTimeout(context.Background(), c.timeout)
	defer cancle()

	return c.client.Get(ctx, key).Result()
}

func (c *Cache) Get(key string, target any) error {
	ctx, cancle := context.WithTimeout(context.Background(), c.timeout)
	defer cancle()

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), target)
}

func (c *Cache) Del(key string) error {
	ctx, cancle := context.WithTimeout(context.Background(), c.timeout)
	defer cancle()

	return c.client.Del(ctx, key).Err()
}
