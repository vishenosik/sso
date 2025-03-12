package app

import (
	"fmt"
	"log/slog"

	cfg "github.com/blacksmith-vish/sso/internal/app/config"
	"github.com/blacksmith-vish/sso/internal/store/cache"
	"github.com/blacksmith-vish/sso/internal/store/cache/providers/noop"
	"github.com/blacksmith-vish/sso/internal/store/cache/providers/redis"
	"github.com/blacksmith-vish/sso/pkg/helpers/config"
	"github.com/blacksmith-vish/sso/pkg/logger/attrs"
)

func (app *App) loadCache() *cache.Cache {

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
		app.log.Error("Failed to init redis cache", attrs.Error(err))
		return cache.NewCache(noop.NewNoopCache())
	}

	app.log.Info(
		"Connected to redis cache",
		slog.String("addr", fmt.Sprintf("%s:%d", conf.Host, conf.Port)),
	)

	return cache.NewCache(redisCache)
}
