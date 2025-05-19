package users

import (
	"context"
	"encoding/json"

	"github.com/dgraph-io/dgo/v240/protos/api"
	"github.com/pkg/errors"
	dmodels "github.com/vishenosik/sso/internal/store/dgraph/models"
	"github.com/vishenosik/sso/internal/store/models"
)

// SaveUser saves user to db.
func (store *store) SaveUser(ctx context.Context, id, nickname, email string, passHash []byte) error {
	const op = "Store.Dgraph.SaveUser"

	variables := map[string]string{
		"$email": email,
	}

	q := `query UserByEmail($email: string){
		users(func: eq(email, $email)) {
			uid
		}
	}`

	txn := store.db.NewTxn()
	defer txn.Discard(ctx)

	resp, err := txn.QueryWithVars(ctx, q, variables)
	if err != nil {
		return errors.Wrap(err, op)
	}

	type Root struct {
		Users []dmodels.User `json:"users"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return errors.Wrap(err, op)
	}

	if len(r.Users) != 0 {
		return errors.Wrap(models.ErrAlreadyExists, op)
	}

	user := &dmodels.User{
		Nickname:     nickname,
		Email:        email,
		ID:           id,
		PasswordHash: passHash,
	}

	userPB, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, op)
	}

	mu := &api.Mutation{
		CommitNow: true,
		SetJson:   userPB,
	}

	// TODO: apply metrics here
	_, err = txn.Mutate(ctx, mu)
	if err != nil {
		// TODO: handle error
		return errors.Wrap(err, op)
	}

	return nil
}
