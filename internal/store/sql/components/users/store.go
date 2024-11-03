package users

import (
	"database/sql"
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
