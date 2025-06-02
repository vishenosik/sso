package sqlite

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/vishenosik/sso/internal/store/models"
)

type AppsStore struct {
	db *sqlx.DB
}

func NewAppsStore(db *sqlx.DB) *AppsStore {
	return &AppsStore{
		db: db,
	}
}

// App returns app by id.
func (store *AppsStore) App(ctx context.Context, id string) (models.App, error) {
	const op = "Store.sqlite.App"

	stmt, err := store.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return models.App{}, errors.Wrap(err, op)
	}

	row := stmt.QueryRowContext(ctx, id)

	var app models.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, errors.Wrap(models.ErrNotFound, op)
		}

		return models.App{}, errors.Wrap(err, op)
	}

	return app, nil
}
