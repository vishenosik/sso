package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
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
