package app

import (
	"context"

	embed "github.com/vishenosik/sso"
	appctx "github.com/vishenosik/sso/internal/app/context"
	"github.com/vishenosik/sso/internal/store/dgraph"
	sqlstore "github.com/vishenosik/sso/internal/store/sql"
	"github.com/vishenosik/sso/internal/store/sql/providers/sqlite"
	"github.com/vishenosik/sso/pkg/helpers/config"
	"github.com/vishenosik/sso/pkg/logger/handlers/std"
	"github.com/vishenosik/sso/pkg/migrate"
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

func loadDgraph(ctx context.Context) (*dgraph.Client, error) {

	appContext := appctx.AppCtx(ctx)
	conf := appContext.Config.Dgraph

	client, err := dgraph.NewClientCtx(
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
		return nil, err
	}

	return client, nil
}
