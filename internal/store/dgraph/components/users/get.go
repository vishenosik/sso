package users

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	dmodels "github.com/vishenosik/sso/internal/store/dgraph/models"
	"github.com/vishenosik/sso/internal/store/models"
)

// SaveUser saves user to db.
func (store *store) UserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "Store.Dgraph.UserByEmail"

	variables := map[string]string{
		"$email": email,
	}

	q := `query UserByEmail($email: string){
		users(func: eq(email, $email)) {
			uuid
			nickname
			email
			pass_hash
		}
	}`

	txn := store.db.NewTxn()
	defer txn.Discard(ctx)

	resp, err := txn.QueryWithVars(ctx, q, variables)
	if err != nil {
		return models.User{}, errors.Wrap(err, op)
	}

	type Root struct {
		Users []dmodels.User `json:"users"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return models.User{}, errors.Wrap(err, op)
	}

	if len(r.Users) == 0 {
		return models.User{}, errors.Wrap(models.ErrNotFound, op)
	}

	user := r.Users[0]

	return models.User{
		Nickname:     user.Nickname,
		Email:        user.Email,
		ID:           user.ID,
		PasswordHash: user.PasswordHash,
	}, nil
}
