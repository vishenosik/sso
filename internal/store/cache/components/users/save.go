package users

import (
	"context"
	"encoding/json"
	"time"

	"github.com/blacksmith-vish/sso/internal/store/models"
)

func (users *Users) SaveUser(ctx context.Context, user models.User) error {
	cacheKey := models.UserCacheKey(user.ID)

	// Clear any cached user data for this ID
	err := users.cache.Delete(ctx, cacheKey)
	if err != nil {
		// TODO: handle error
		return err
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		// TODO: handle error
		return err
	}

	// Save the user to the cache
	err = users.cache.Set(ctx, cacheKey, string(userJSON), time.Hour)
	if err != nil {
		// TODO: handle error
		return err
	}

	return nil
}
