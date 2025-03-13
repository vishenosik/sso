package app

import (
	"context"
	"log/slog"
	"time"

	embed "github.com/blacksmith-vish/sso"
	cfg "github.com/blacksmith-vish/sso/internal/app/config"
	grpcApp "github.com/blacksmith-vish/sso/internal/app/grpc"
	restApp "github.com/blacksmith-vish/sso/internal/app/rest"
	authenticationService "github.com/blacksmith-vish/sso/internal/services/authentication"
	"github.com/blacksmith-vish/sso/internal/store/combined"
	"github.com/blacksmith-vish/sso/internal/store/dgraph"
	sqlstore "github.com/blacksmith-vish/sso/internal/store/sql"

	libctx "github.com/blacksmith-vish/sso/internal/lib/context"
	"github.com/blacksmith-vish/sso/internal/store/sql/providers/sqlite"
	"github.com/blacksmith-vish/sso/pkg/helpers/config"
	"github.com/blacksmith-vish/sso/pkg/logger/attrs"
	"github.com/blacksmith-vish/sso/pkg/logger/handlers/std"
	"github.com/blacksmith-vish/sso/pkg/migrate"
)

const (
	envDev  = "dev"
	envProd = "prod"
	envTest = "test"
)

type App struct {
	grpcServer *grpcApp.App
	restServer *restApp.App
	log        *slog.Logger
}

func NewApp() (*App, error) {

	conf := cfg.EnvConfig()

	// logger setup
	// TODO: implement env logic
	log := setupLogger(conf.Env)

	appctx, err := libctx.NewAppContext(log)
	if err != nil {
		return nil, err
	}

	ctx := libctx.SetAppContext(context.Background(), appctx)

	app := &App{log: log}

	// Cache init
	cache := app.loadCache()

	// Stores init
	sqliteStore := sqlite.MustInitSqlite(conf.StorePath)
	store := sqlstore.NewStore(sqliteStore)

	_, err = dgraph.NewClient(ctx, dgraph.Config{
		Credentials: config.Credentials{
			User:     conf.Dgraph.User,
			Password: conf.Dgraph.Password,
		},
		GrpcServer: config.Server{
			Host: conf.Dgraph.GrpcHost,
			Port: conf.Dgraph.GrpcPort,
		},
	})

	log.Debug("", slog.Any("conf", conf))

	if err != nil {
		log.Error(
			"failed to connect to dgraph",
			attrs.Error(err),
		)
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

	app.grpcServer = grpcApp.NewGrpcApp(
		log,
		grpcApp.Config{
			Server: config.Server{
				Port: conf.GrpcConfig.Port,
			},
		},
		authenticationService,
	)

	app.restServer = restApp.NewRestApp(
		ctx,
		restApp.Config{
			Server: config.Server{
				Port: conf.RestConfig.Port,
			},
		},
		authenticationService,
	)

	return app, nil
}

func (app *App) MustRun() {

	app.log.Info("start app")

	// Инициализация gRPC-сервер
	go app.grpcServer.MustRun()

	go app.restServer.MustRun()
}

func (app *App) Stop(signal string) {

	app.log.Info("app stopping", slog.String("signal", signal))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer func() {
		// extra handling here
		cancel()
	}()

	app.grpcServer.Stop()

	app.restServer.Stop(ctx)

	app.log.Info("app stopped")
}
