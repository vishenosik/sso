package app

import (
	"context"
	"fmt"
	"log/slog"

	appctx "github.com/vishenosik/sso/internal/app/context"
	"github.com/vishenosik/sso/internal/store/cache"
	"github.com/vishenosik/sso/internal/store/cache/providers/noop"
	"github.com/vishenosik/sso/internal/store/cache/providers/redis"
	"github.com/vishenosik/sso/pkg/helpers/config"
	"github.com/vishenosik/sso/pkg/logger/attrs"
)

func loadCache(ctx context.Context) *cache.Cache {

	appContext := appctx.AppCtx(ctx)
	log := appContext.Logger
	conf := appContext.Config.Redis

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
