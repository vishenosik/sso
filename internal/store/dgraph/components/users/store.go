package users

import (
	"github.com/dgraph-io/dgo/v210"
)

type Users struct {
	db *dgo.Dgraph
}

type store = Users

func NewUsersStore(db *dgo.Dgraph) *Users {
	return &Users{
		db: db,
	}
}
