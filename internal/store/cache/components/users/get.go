package users

import (
	"context"
	"encoding/json"

	"github.com/blacksmith-vish/sso/internal/store/models"
)

func (users *Users) UserByID(ctx context.Context, id string) (models.User, error) {
	var user models.User
	cacheKey := models.UserCacheKey(id)

	// Try to get user from cache
	jsonUser, err := users.cache.Get(ctx, cacheKey)
	if err != nil {
		// TODO: handle error
		return models.User{}, err
	}

	if err := json.Unmarshal([]byte(jsonUser), &user); err != nil {
		// TODO: handle error
		return models.User{}, err
	}

	return user, nil
}
