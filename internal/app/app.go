package app

import (
	"log/slog"

	authentication "github.com/blacksmith-vish/sso/internal/api/authentication/rest"
	grpcApp "github.com/blacksmith-vish/sso/internal/app/grpc"
	restApp "github.com/blacksmith-vish/sso/internal/app/rest"
	"github.com/blacksmith-vish/sso/internal/lib/config"
	"github.com/blacksmith-vish/sso/internal/lib/migrate"
	authService "github.com/blacksmith-vish/sso/internal/services/authentication"
	sqlstore "github.com/blacksmith-vish/sso/internal/store/sql"
	"github.com/blacksmith-vish/sso/internal/store/sql/sqlite"
)

type App struct {
	GRPCServer *grpcApp.App
	RESTServer *restApp.App
}

func NewApp(
	log *slog.Logger,
	conf *config.Config,
) *App {

	// Инициализация хранилища
	sqliteStore := sqlite.MustInitSqlite(conf.StorePath)
	migrate.MustMigrate(sqliteStore)

	store := sqlstore.NewStore(sqliteStore)

	// Инициализация auth сервиса
	authService := authService.NewService(
		log,
		conf.AuthenticationService,
		store.AuthenticationStore(),
		store.AuthenticationStore(),
		store.AuthenticationStore(),
	)

	grpcapp := grpcApp.NewGrpcApp(log, conf.GrpcConfig, authService)

	restapp := restApp.NewRestApp(log, conf.RestConfig, authentication.NewAuthenticationServer(log, authService))

	return &App{
		GRPCServer: grpcapp,
		RESTServer: restapp,
	}

}
