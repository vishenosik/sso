package sqlstore

import (
	"database/sql"

	"github.com/blacksmith-vish/sso/internal/store/sql/components/apps"
	"github.com/blacksmith-vish/sso/internal/store/sql/components/users"
)

type Store struct {
	db    *sql.DB
	apps  *apps.Apps
	users *users.Users
}

type StoreProvider interface {
	DB() *sql.DB
}

func NewStore(
	store StoreProvider,
) *Store {

	db := store.DB()

	return &Store{
		db:    db,
		apps:  apps.NewAppsStore(db),
		users: users.NewUsersStore(db),
	}
}

func (store *Store) Stop() error {
	return store.db.Close()
}

func (store *Store) Apps() *apps.Apps {
	return store.apps
}

func (store *Store) Users() *users.Users {
	return store.users
}
