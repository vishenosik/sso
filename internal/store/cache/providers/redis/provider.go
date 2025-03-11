package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/blacksmith-vish/sso/pkg/helpers/config"
	"github.com/go-redis/redis/v8"
)

type Config struct {
	Server      config.Server
	Credentials config.Credentials
	DB          int
}

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(config Config) (*redisCache, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port),
		Username: config.Credentials.User,
		Password: config.Credentials.Password,
		DB:       config.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &redisCache{client: client}, nil
}

func (ca *redisCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return ca.client.Set(ctx, key, value, expiration).Err()
}

func (ca *redisCache) Get(ctx context.Context, key string) (string, error) {
	return ca.client.Get(ctx, key).Result()
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
