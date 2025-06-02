package users

import (
	"context"
	"database/sql"

	"github.com/vishenosik/sso/internal/store/models"

	"github.com/pkg/errors"
)

// User returns user by email.
func (store *store) UserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "Store.sqlite.User"

	stmt, err := store.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
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
