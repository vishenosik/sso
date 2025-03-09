package app

import (
	"context"
	"log/slog"
	"time"

	embed "github.com/blacksmith-vish/sso"
	"github.com/blacksmith-vish/sso/internal/app/config"
	grpcApp "github.com/blacksmith-vish/sso/internal/app/grpc"
	restApp "github.com/blacksmith-vish/sso/internal/app/rest"
	cfg "github.com/blacksmith-vish/sso/internal/lib/config"
	authenticationService "github.com/blacksmith-vish/sso/internal/services/authentication"
	"github.com/blacksmith-vish/sso/internal/store/combined"
	config_yaml "github.com/blacksmith-vish/sso/internal/store/filesystem/config/yaml"
	sqlstore "github.com/blacksmith-vish/sso/internal/store/sql"
	"github.com/blacksmith-vish/sso/internal/store/sql/providers/sqlite"
	"github.com/blacksmith-vish/sso/pkg/logger"
	"github.com/blacksmith-vish/sso/pkg/logger/handlers/std"
	"github.com/blacksmith-vish/sso/pkg/migrate"
)

type App struct {
	grpcServer *grpcApp.App
	restServer *restApp.App
	config     config.ConfigSources
	log        *slog.Logger
}

func NewApp() *App {

	// Инициализация конфига
	conf := cfg.NewConfig(config_yaml.MustLoad())

	app := &App{
		log:    logger.SetupLogger(conf.Env),
		config: config.NewConfigSources(),
	}

	// Stores init
	sqliteStore := sqlite.MustInitSqlite(conf.StorePath)
	store := sqlstore.NewStore(sqliteStore)

	// Stores migration
	migrate.NewMigrator(
		std.NewStdLogger(app.log),
		embed.SQLiteMigrations,
	).MustMigrate(sqliteStore)

	// Cache init
	cache := app.redisCache(app.log)

	// Data schemas init
	cachedStore := combined.NewCachedStore(store, cache)

	// Services init
	authenticationService := authenticationService.NewService(
		app.log,
		conf.AuthenticationService,
		store,
		store,
		cachedStore,
	)

	// GRPC services init
	app.grpcServer = grpcApp.NewGrpcApp(app.log, conf.GrpcConfig, authenticationService)

	// REST services init
	app.restServer = restApp.NewRestApp(app.log, conf.RestConfig, authenticationService)

	return app
}

func (app *App) Run() {

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
