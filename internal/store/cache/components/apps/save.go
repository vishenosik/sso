package apps

import (
	"context"
	"encoding/json"
	"time"

	"github.com/blacksmith-vish/sso/internal/store/models"
)

func (apps *Apps) SaveApp(ctx context.Context, app models.App) error {
	cacheKey := models.AppCacheKey(app.ID)

	// Clear any cached user data for this ID
	err := apps.cache.Delete(ctx, cacheKey)
	if err != nil {
		// TODO: handle error
		return err
	}

	appJSON, err := json.Marshal(app)
	if err != nil {
		// TODO: handle error
		return err
	}

	// Save the user to the cache
	err = apps.cache.Set(ctx, cacheKey, string(appJSON), time.Hour*24)
	if err != nil {
		// TODO: handle error
		return err
	}

	return nil
}
