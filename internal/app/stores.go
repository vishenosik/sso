package app

import (
	"context"

	embed "github.com/blacksmith-vish/sso"
	appctx "github.com/blacksmith-vish/sso/internal/app/context"
	"github.com/blacksmith-vish/sso/internal/store/dgraph"
	sqlstore "github.com/blacksmith-vish/sso/internal/store/sql"
	"github.com/blacksmith-vish/sso/internal/store/sql/providers/sqlite"
	"github.com/blacksmith-vish/sso/pkg/helpers/config"
	"github.com/blacksmith-vish/sso/pkg/logger/handlers/std"
	"github.com/blacksmith-vish/sso/pkg/migrate"
	"github.com/dgraph-io/dgo/v210"
)

func loadSqlStore(ctx context.Context) (*sqlstore.Store, error) {

	appContext := appctx.AppCtx(ctx)

	// Stores init
	sqliteStore := sqlite.MustInitSqlite(appContext.Config.StorePath)
	store := sqlstore.NewStore(sqliteStore)

	// Stores migration
	migrate.NewMigrator(
		std.NewStdLogger(appContext.Logger),
		embed.SQLiteMigrations,
	).MustMigrate(sqliteStore)

	return store, nil
}

func loadDgraph(ctx context.Context) (*dgo.Dgraph, error) {

	appContext := appctx.AppCtx(ctx)
	conf := appContext.Config.Dgraph

	client, err := dgraph.NewClient(
		ctx,
		dgraph.Config{
			Credentials: config.Credentials{
				User:     conf.User,
				Password: conf.Password,
			},
			GrpcServer: config.Server{
				Host: conf.GrpcHost,
				Port: conf.GrpcPort,
			},
		},
	)

	if err != nil {
		// return nil, err
	}

	return client, nil
}
