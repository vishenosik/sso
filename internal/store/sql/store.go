package sqlstore

import (
	"context"
	"database/sql"

	"github.com/blacksmith-vish/sso/internal/store/models"
	"github.com/blacksmith-vish/sso/internal/store/sql/components/apps"
	"github.com/blacksmith-vish/sso/internal/store/sql/components/users"
)

type Store struct {
	provider StoreProvider
	apps     *apps.Apps
	users    *users.Users
}

type StoreProvider interface {
	DB() *sql.DB
}

func NewStore(
	provider StoreProvider,
) *Store {

	db := provider.DB()

	return &Store{
		provider: provider,
		apps:     apps.NewAppsStore(db),
		users:    users.NewUsersStore(db),
	}
}

func (store *Store) Stop() error {
	return store.provider.DB().Close()
}

func (store *Store) App(ctx context.Context, id string) (models.App, error) {
	return store.apps.App(ctx, id)
}

func (store *Store) IsAdmin(ctx context.Context, userID string) (bool, error) {
	return store.users.IsAdmin(ctx, userID)
}

func (store *Store) SaveUser(ctx context.Context, id string, nickname string, email string, passHash []byte) error {
	return store.users.SaveUser(ctx, id, nickname, email, passHash)
}

func (store *Store) UserByEmail(ctx context.Context, email string) (models.User, error) {
	return store.users.UserByEmail(ctx, email)
}
