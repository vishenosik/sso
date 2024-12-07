package app

import (
	"github.com/blacksmith-vish/sso/internal/lib/migrate"
	sqlstore "github.com/blacksmith-vish/sso/internal/store/sql"
	"github.com/blacksmith-vish/sso/internal/store/sql/providers/sqlite"
)

func sqliteStore(StorePath string) *sqlstore.Store {
	sqliteStore := sqlite.MustInitSqlite(StorePath)
	migrate.MustMigrate(sqliteStore)
	return sqlstore.NewStore(sqliteStore)
}
