package app

import (
	cfg "github.com/blacksmith-vish/sso/internal/lib/config"
	"github.com/blacksmith-vish/sso/internal/store/cache"
	"github.com/blacksmith-vish/sso/internal/store/cache/providers/redis"
	"github.com/blacksmith-vish/sso/pkg/helpers/config"
)

func loadCache() *cache.Cache {

	conf := cfg.EnvConfig().Redis

	redisCache, err := redis.NewRedisCache(redis.Config{
		Credentials: config.Credentials{
			User:     conf.User,
			Password: conf.Password,
		},
		Server: config.Server{
			Host: conf.Host,
			Port: conf.Port,
		},
		DB: conf.DB,
	})
	if err != nil {
		// TODO: handle error
		// log.Error("Failed to init redis cache", attrs.Error(err))
		// return cache.NewCache(noop.NewNoopCache())
	}

	return cache.NewCache(redisCache)
}
