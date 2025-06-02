package apps

import (
	"database/sql"
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
