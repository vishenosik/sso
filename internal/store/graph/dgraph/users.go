package dgraph

import (
	"context"
	"encoding/json"

	"github.com/dgraph-io/dgo/v240"
	"github.com/dgraph-io/dgo/v240/protos/api"
	"github.com/vishenosik/gocherry/pkg/errors"
	"github.com/vishenosik/sso/internal/entities"
	"github.com/vishenosik/sso/internal/store/graph/models"
)

type UsersStore struct {
	db *dgo.Dgraph
}

func NewUsersStore(db *dgo.Dgraph) *UsersStore {
	return &UsersStore{
		db: db,
	}
}

// SaveUser saves user to db.
func (us *UsersStore) SaveUser(ctx context.Context, user *entities.UserCreds) error {
	return us.save(ctx, models.UserFromEntities(user))
}

func (us *UsersStore) save(ctx context.Context, user *models.User) error {
	const op = "Store.Dgraph.SaveUser"

	variables := map[string]string{
		"$email": user.Email,
	}

	q := `query UserByEmail($email: string){
		users(func: eq(email, $email)) {
			uid
		}
	}`

	txn := us.db.NewTxn()
	defer txn.Discard(ctx)

	resp, err := txn.QueryWithVars(ctx, q, variables)
	if err != nil {
		return errors.Wrap(err, op)
	}

	type Root struct {
		Users []models.User `json:"users"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return errors.Wrap(err, op)
	}

	if len(r.Users) != 0 {
		return errors.Wrap(entities.ErrAlreadyExists, op)
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

// SaveUser saves user to db.
func (us *UsersStore) UserByEmail(ctx context.Context, email string) (*entities.UserCreds, error) {
	user, err := us.userByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return models.UserToEntities(user), nil
}

// SaveUser saves user to db.
func (us *UsersStore) userByEmail(ctx context.Context, email string) (*models.User, error) {
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

	txn := us.db.NewTxn()
	defer txn.Discard(ctx)

	resp, err := txn.QueryWithVars(ctx, q, variables)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	type Root struct {
		Users []models.User `json:"users"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	if len(r.Users) == 0 {
		return nil, errors.Wrap(entities.ErrNotFound, op)
	}

	user := r.Users[0]

	return &models.User{
		Nickname:     user.Nickname,
		Email:        user.Email,
		ID:           user.ID,
		PasswordHash: user.PasswordHash,
	}, nil
}

func (store *UsersStore) IsAdmin(ctx context.Context, userID string) (isAdmin bool, err error) {
	return false, nil
}
