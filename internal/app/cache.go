package app

import (
	"log/slog"

	"github.com/blacksmith-vish/sso/internal/store/cache"
	"github.com/blacksmith-vish/sso/internal/store/cache/providers/noop"
	"github.com/blacksmith-vish/sso/internal/store/cache/providers/redis"
	"github.com/blacksmith-vish/sso/pkg/logger/attrs"
)

func (app *App) redisCache(
	log *slog.Logger,
) *cache.Cache {
	redisCache, err := redis.NewRedisCache(app.config.FS)
	if err != nil {
		// TODO: handle error
		log.Error("Failed to init redis cache", attrs.Error(err))
		return cache.NewCache(noop.NewNoopCache())
	}
	return cache.NewCache(redisCache)
}
