package apps

import (
	"context"
	"encoding/json"

	"github.com/blacksmith-vish/sso/internal/store/models"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

func (apps *Apps) App(ctx context.Context, id string) (models.App, error) {
	const op = "store.cache.app"

	jsonApp, err := apps.cache.Get(ctx, models.AppCacheKey(id))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return models.App{}, errors.Wrap(models.ErrNotFound, op)
		}
		return models.App{}, errors.Wrap(err, op)
	}

	var app models.App
	if err := json.Unmarshal([]byte(jsonApp), &app); err != nil {
		return models.App{}, errors.Wrap(err, op)
	}
	return app, nil
}
