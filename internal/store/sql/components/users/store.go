package users

import (
	"database/sql"

	"github.com/pkg/errors"
)

var (
	ErrUserExists = errors.New("user exists already")
	// user not found
	ErrUserNotFound = errors.New("user not found")
)

type Users struct {
	db *sql.DB
}

type store = Users

func NewUsersStore(db *sql.DB) *Users {
	return &Users{
		db: db,
	}
}
