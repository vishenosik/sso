package app

import (
	"context"
	"fmt"
	"log/slog"

	cfg "github.com/blacksmith-vish/sso/internal/app/config"
	"github.com/blacksmith-vish/sso/internal/store/cache"
	"github.com/blacksmith-vish/sso/internal/store/cache/providers/noop"
	"github.com/blacksmith-vish/sso/internal/store/cache/providers/redis"
	"github.com/blacksmith-vish/sso/pkg/helpers/config"
	"github.com/blacksmith-vish/sso/pkg/logger/attrs"

	libctx "github.com/blacksmith-vish/sso/internal/lib/context"
)

func loadCache(ctx context.Context) *cache.Cache {

	appctx, ok := libctx.AppCtx(ctx)
	if !ok {
		// TODO: handle error
	}

	log := appctx.Logger

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
		log.Error("Failed to init redis cache", attrs.Error(err))
		return cache.NewCache(noop.NewNoopCache())
	}

	log.Info(
		"Connected to redis cache",
		slog.String("addr", fmt.Sprintf("%s:%d", conf.Host, conf.Port)),
	)

	return cache.NewCache(redisCache)
}
