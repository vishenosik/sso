package app

import (
	"log/slog"

	grpcApp "github.com/blacksmith-vish/sso/internal/app/grpc"
	restApp "github.com/blacksmith-vish/sso/internal/app/rest"
	"github.com/blacksmith-vish/sso/internal/lib/config"
	"github.com/blacksmith-vish/sso/internal/lib/migrate"
	authenticationService "github.com/blacksmith-vish/sso/internal/services/authentication"
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

	// Services init
	authenticationService := authenticationService.NewService(
		log,
		conf.AuthenticationService,
		store.Users(),
		store.Users(),
		store.Apps(),
	)

	grpcapp := grpcApp.NewGrpcApp(log, conf.GrpcConfig, authenticationService)

	restapp := restApp.NewRestApp(log, conf.RestConfig, authenticationService)

	return &App{
		GRPCServer: grpcapp,
		RESTServer: restapp,
	}

}
