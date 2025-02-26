package app

import (
	"log/slog"

	embed "github.com/blacksmith-vish/sso"
	grpcApp "github.com/blacksmith-vish/sso/internal/app/grpc"
	restApp "github.com/blacksmith-vish/sso/internal/app/rest"
	"github.com/blacksmith-vish/sso/internal/lib/config"
	authenticationService "github.com/blacksmith-vish/sso/internal/services/authentication"
	"github.com/blacksmith-vish/sso/internal/store/combined"
	sqlstore "github.com/blacksmith-vish/sso/internal/store/sql"
	"github.com/blacksmith-vish/sso/internal/store/sql/providers/sqlite"
	"github.com/blacksmith-vish/sso/pkg/logger/handlers/std"
	"github.com/blacksmith-vish/sso/pkg/migrate"
)

type App struct {
	GRPCServer *grpcApp.App
	RESTServer *restApp.App
}

func NewApp(
	log *slog.Logger,
	conf *config.Config,
) *App {

	// Stores init
	sqliteStore := sqlite.MustInitSqlite(conf.StorePath)
	store := sqlstore.NewStore(sqliteStore)

	// Stores migration
	migrate.NewMigrator(
		std.NewStdLogger(log),
		embed.SQLiteMigrations,
	).MustMigrate(sqliteStore)

	// Cache init
	cache := redisCache(log, conf.Redis)

	// Data schemas init
	cachedStore := combined.NewCachedStore(store, cache)

	// Services init
	authenticationService := authenticationService.NewService(
		log,
		conf.AuthenticationService,
		store,
		store,
		cachedStore,
	)

	// GRPC services init
	grpcapp := grpcApp.NewGrpcApp(log, conf.GrpcConfig, authenticationService)

	// REST services init
	restapp := restApp.NewRestApp(log, conf.RestConfig, authenticationService)

	return &App{
		GRPCServer: grpcapp,
		RESTServer: restapp,
	}
}
