package app

import (
	"log/slog"

	grpcApp "github.com/blacksmith-vish/sso/internal/app/grpc"
	restApp "github.com/blacksmith-vish/sso/internal/app/rest"
	"github.com/blacksmith-vish/sso/internal/lib/config"
	"github.com/blacksmith-vish/sso/internal/lib/logger/attrs"
	"github.com/blacksmith-vish/sso/internal/lib/migrate"
	authenticationService "github.com/blacksmith-vish/sso/internal/services/authentication"
	"github.com/blacksmith-vish/sso/internal/store/cache"
	"github.com/blacksmith-vish/sso/internal/store/cache/providers/redis"
	"github.com/blacksmith-vish/sso/internal/store/combined"
	sqlstore "github.com/blacksmith-vish/sso/internal/store/sql"
	"github.com/blacksmith-vish/sso/internal/store/sql/providers/sqlite"
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
	migrate.MustMigrate(sqliteStore)

	store := sqlstore.NewStore(sqliteStore)

	// Cache init
	redisCache, err := redis.NewRedisCache("")
	if err != nil {
		// TODO: handle error
		log.Error("Failed to create redis cache", attrs.Error(err))
	}
	cache := cache.NewCache(redisCache)

	// Data schemas init

	cachedDB := combined.NewCachedDB(store, cache)

	// Services init
	authenticationService := authenticationService.NewService(
		log,
		conf.AuthenticationService,
		store,
		store,
		cachedDB,
	)

	grpcapp := grpcApp.NewGrpcApp(log, conf.GrpcConfig, authenticationService)

	restapp := restApp.NewRestApp(log, conf.RestConfig, authenticationService)

	return &App{
		GRPCServer: grpcapp,
		RESTServer: restapp,
	}

}
