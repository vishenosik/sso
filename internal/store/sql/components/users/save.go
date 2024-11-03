package users

import (
	"context"

	"github.com/blacksmith-vish/sso/internal/store/models"
	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

// SaveUser saves user to db.
func (store *store) SaveUser(ctx context.Context, id, nickname, email string, passHash []byte) error {
	const op = "Store.sqlite.SaveUser"

	stmt, err := store.db.Prepare("INSERT INTO users(id, nickname, email, pass_hash) VALUES(?, ?, ?, ?)")
	if err != nil {
		return errors.Wrap(err, op)
	}

	_, err = stmt.ExecContext(ctx, id, nickname, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return errors.Wrap(models.ErrAlreadyExists, op)
		}

		return errors.Wrap(err, op)
	}

	return nil
}
