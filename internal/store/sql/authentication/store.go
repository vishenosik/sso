package authentication

import (
	"database/sql"

	"github.com/pkg/errors"
)

const (
	MethodApp      = "App"
	methodIsAdmin  = "IsAdmin"
	methodSaveUser = "SaveUser"
	methodUser     = "User"
)

var (
	ErrAppNotFound = errors.New("app not found")
	ErrUserExists  = errors.New("user exists already")
	// user not found
	ErrUserNotFound = errors.New("user not found")
)

type AuthenticationStore struct {
	db *sql.DB
}

type Store = AuthenticationStore

func NewAuthenticationStore(db *sql.DB) *AuthenticationStore {
	return &AuthenticationStore{
		db: db,
	}
}
