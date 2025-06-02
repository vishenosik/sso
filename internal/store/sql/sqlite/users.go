package sqlite

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/vishenosik/sso/internal/entities"
	"github.com/vishenosik/sso/internal/store/sql/models"
)

type UsersStore struct {
	db *sqlx.DB
}

func NewUsersStore(db *sqlx.DB) *UsersStore {
	return &UsersStore{
		db: db,
	}
}

// SaveUser saves user to db.
func (store *UsersStore) SaveUser(ctx context.Context, user *entities.UserCreds) error {
	const op = "Store.sqlite.SaveUser"

	stmt, err := store.db.Prepare("INSERT INTO UsersStore(id, nickname, email, pass_hash) VALUES(?, ?, ?, ?)")
	if err != nil {
		return errors.Wrap(err, op)
	}

	_, err = stmt.ExecContext(ctx, user.ID, user.Nickname, user.Email, user.PasswordHash)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

// User returns user by email.
func (store *UsersStore) UserByEmail(ctx context.Context, email string) (*entities.UserCreds, error) {
	const op = "Store.sqlite.User"

	stmt, err := store.db.Prepare("SELECT id, email, pass_hash FROM UsersStore WHERE email = ?")
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	row := stmt.QueryRowContext(ctx, email)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(entities.ErrNotFound, op)
		}

		return nil, errors.Wrap(err, op)
	}

	return &entities.UserCreds{
		User: entities.User{
			ID:    user.ID,
			Email: user.Email,
		},
		Password: string(user.Password),
	}, nil
}

func (store *UsersStore) IsAdmin(ctx context.Context, userID string) (isAdmin bool, err error) {
	return false, nil
}
