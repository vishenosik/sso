package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	embed "github.com/blacksmith-vish/sso"
	grpcApp "github.com/blacksmith-vish/sso/internal/app/grpc"
	restApp "github.com/blacksmith-vish/sso/internal/app/rest"
	authenticationService "github.com/blacksmith-vish/sso/internal/services/authentication"
	"github.com/blacksmith-vish/sso/internal/store/combined"
	"github.com/blacksmith-vish/sso/internal/store/dgraph"
	sqlstore "github.com/blacksmith-vish/sso/internal/store/sql"

	appctx "github.com/blacksmith-vish/sso/internal/app/context"
	"github.com/blacksmith-vish/sso/internal/store/sql/providers/sqlite"
	"github.com/blacksmith-vish/sso/pkg/helpers/config"
	"github.com/blacksmith-vish/sso/pkg/logger/handlers/std"
	"github.com/blacksmith-vish/sso/pkg/migrate"
)

type App struct {
	grpcServer *grpcApp.App
	restServer *restApp.App
	log        *slog.Logger
}

func MustInitApp() *App {
	app, err := NewApp()
	if err != nil {
		panic(fmt.Sprintf("failed to create app %s", err))
	}
	return app
}

func NewApp() (*App, error) {

	ctx := appctx.SetupAppCtx()

	appContext, ok := appctx.AppCtx(ctx)
	if !ok {
		return nil, errors.New("failed to get app context")
	}

	// logger setup
	// TODO: implement env logic
	log := appContext.Logger
	conf := appContext.Config

	// Cache init
	cache := loadCache(ctx)

	// Stores init
	sqliteStore := sqlite.MustInitSqlite(conf.StorePath)
	store := sqlstore.NewStore(sqliteStore)

	_, err := dgraph.NewClient(
		ctx,
		dgraph.Config{
			Credentials: config.Credentials{
				User:     conf.Dgraph.User,
				Password: conf.Dgraph.Password,
			},
			GrpcServer: config.Server{
				Host: conf.Dgraph.GrpcHost,
				Port: conf.Dgraph.GrpcPort,
			},
		},
	)

	if err != nil {
		// return nil, err
	}

	// Stores migration
	migrate.NewMigrator(
		std.NewStdLogger(log),
		embed.SQLiteMigrations,
	).MustMigrate(sqliteStore)

	// Data schemas init
	cachedStore := combined.NewCachedStore(store, cache)

	// Services init
	authenticationService := authenticationService.NewService(
		ctx,
		log,
		conf.AuthenticationService,
		store,
		store,
		cachedStore,
	)

	return &App{
		log: log,
		grpcServer: grpcApp.NewGrpcApp(
			log,
			grpcApp.Config{
				Server: config.Server{
					Port: conf.GrpcConfig.Port,
				},
			},
			authenticationService,
		),
		restServer: restApp.NewRestApp(
			ctx,
			restApp.Config{
				Server: config.Server{
					Port: conf.RestConfig.Port,
				},
			},
			authenticationService,
		),
	}, nil
}

func (app *App) MustRun() {

	app.log.Info("start app")

	// Инициализация gRPC-сервер
	go app.grpcServer.MustRun()

	go app.restServer.MustRun()
}

func (app *App) Stop(ctx context.Context) {

	const msg = "app stopping"

	signal, ok := appctx.SignalCtx(ctx)
	if !ok {
		app.log.Info(msg, slog.String("signal", signal.Signal.String()))
	} else {
		app.log.Info(msg)
	}

	app.grpcServer.Stop()

	app.restServer.Stop(ctx)

	app.log.Info("app stopped")
}
