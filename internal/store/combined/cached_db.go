package combined

import (
	"context"

	"github.com/vishenosik/sso/internal/store/cache"
	"github.com/vishenosik/sso/internal/store/models"
	sqlstore "github.com/vishenosik/sso/internal/store/sql"
)

type CachedStore struct {
	store *sqlstore.Store
	cache *cache.Cache
}

func NewCachedStore(
	store *sqlstore.Store,
	cache *cache.Cache,
) *CachedStore {

	if store == nil {
		panic("store cannot be nil")
	}

	if cache == nil {
		// TODO: handle cache initialization error logging
	}

	return &CachedStore{
		store: store,
		cache: cache,
	}
}

func (cdb *CachedStore) App(ctx context.Context, id string) (models.App, error) {

	app, err := cdb.cache.App(ctx, id)
	if err == nil {
		return app, nil
	}

	// TODO: handle error logging

	app, err = cdb.store.App(ctx, id)
	if err != nil {
		// TODO: handle error
		return models.App{}, err
	}

	if err := cdb.cache.SaveApp(ctx, app); err != nil {
		// TODO: handle error
	}

	return app, nil
}

func (cdb *CachedStore) IsAdmin(ctx context.Context, userID string) (bool, error) {
	return cdb.store.IsAdmin(ctx, userID)
}

func (cdb *CachedStore) SaveUser(ctx context.Context, id string, nickname string, email string, passHash []byte) error {
	return cdb.store.SaveUser(ctx, id, nickname, email, passHash)
}

func (cdb *CachedStore) UserByEmail(ctx context.Context, email string) (models.User, error) {
	return cdb.store.UserByEmail(ctx, email)
}
