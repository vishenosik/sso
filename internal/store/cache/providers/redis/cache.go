package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/blacksmith-vish/sso/internal/lib/config"
	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(conf config.Redis) (*RedisCache, error) {

	fmt.Println("configuring Redis cache provider...", conf)

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Options.Host, conf.Options.Port),
		Username: conf.Options.User,
		Password: conf.Options.Password,
		DB:       conf.Options.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisCache{client: client}, nil
}

func (ca *RedisCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return ca.client.Set(ctx, key, value, expiration).Err()
}

func (ca *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return ca.client.Get(ctx, key).Result()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
