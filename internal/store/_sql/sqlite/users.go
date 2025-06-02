package sqlite

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/vishenosik/sso/internal/store/models"
)

type UsersStore struct {
	db *sqlx.DB
}

func NewUsersStoreStore(db *sqlx.DB) *UsersStore {
	return &UsersStore{
		db: db,
	}
}

// SaveUser saves user to db.
func (store *UsersStore) SaveUser(ctx context.Context, id, nickname, email string, passHash []byte) error {
	const op = "Store.sqlite.SaveUser"

	stmt, err := store.db.Prepare("INSERT INTO UsersStore(id, nickname, email, pass_hash) VALUES(?, ?, ?, ?)")
	if err != nil {
		return errors.Wrap(err, op)
	}

	_, err = stmt.ExecContext(ctx, id, nickname, email, passHash)
	if err != nil {
		// var sqliteErr sqlite3.Error
		// if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		// 	return errors.Wrap(models.ErrAlreadyExists, op)
		// }

		return errors.Wrap(err, op)
	}

	return nil
}

func (store *UsersStore) IsAdmin(ctx context.Context, userID string) (bool, error) {
	const op = "Store.sqlite.IsAdmin"

	stmt, err := store.db.Prepare("SELECT is_admin FROM UsersStore WHERE id = ?")
	if err != nil {
		return false, errors.Wrap(err, op)
	}

	row := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool

	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errors.Wrap(models.ErrNotFound, op)
		}

		return false, errors.Wrap(err, op)
	}

	return isAdmin, nil
}

// User returns user by email.
func (store *UsersStore) UserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "Store.sqlite.User"

	stmt, err := store.db.Prepare("SELECT id, email, pass_hash FROM UsersStore WHERE email = ?")
	if err != nil {
		return models.User{}, errors.Wrap(err, op)
	}

	row := stmt.QueryRowContext(ctx, email)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errors.Wrap(models.ErrNotFound, op)
		}

		return models.User{}, errors.Wrap(err, op)
	}

	return user, nil
}
