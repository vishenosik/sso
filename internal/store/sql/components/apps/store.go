package apps

import (
	"database/sql"

	"github.com/pkg/errors"
)

var (
	ErrAppNotFound = errors.New("app not found")
)

type Apps struct {
	db *sql.DB
}

type store = Apps

func NewAppsStore(db *sql.DB) *Apps {
	return &Apps{
		db: db,
	}
}
