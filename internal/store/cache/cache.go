package cache

import (
	"context"

	"github.com/vishenosik/sso/internal/store/cache/components/apps"
	"github.com/vishenosik/sso/internal/store/cache/components/users"
	"github.com/vishenosik/sso/internal/store/cache/providers/noop"
	"github.com/vishenosik/sso/internal/store/models"
)

type Cache struct {
	provider models.CacheProvider
	apps     *apps.Apps
	users    *users.Users
}

func NewCache(
	provider models.CacheProvider,
) *Cache {
	if provider == nil {
		// TODO: handle provider nil
		provider = noop.NewNoopCache()
	}

	return &Cache{
		provider: provider,
		apps:     apps.NewAppsCache(provider),
		users:    users.NewUsersCache(provider),
	}
}

func (ca *Cache) App(ctx context.Context, id string) (models.App, error) {
	return ca.apps.App(ctx, id)
}

func (ca *Cache) SaveApp(ctx context.Context, app models.App) error {
	return ca.apps.SaveApp(ctx, app)
}

func (ca *Cache) UserByID(ctx context.Context, userID string) (models.User, error) {
	return ca.users.UserByID(ctx, userID)
}

func (ca *Cache) SaveUser(ctx context.Context, id string, nickname string, email string, passHash []byte) error {
	return ca.users.SaveUser(
		ctx,
		models.User{
			Nickname: nickname,
			Email:    email,
			ID:       id,
		},
	)
}
