package app

import (
	// std
	"context"
	"net/http"
	"path"

	// pkg
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/vishenosik/gocherry/pkg/grpc"

	// internal pkg
	"github.com/vishenosik/gocherry"
	"github.com/vishenosik/gocherry/pkg/cache"
	_http "github.com/vishenosik/gocherry/pkg/http"
	"github.com/vishenosik/gocherry/pkg/sql"
	"github.com/vishenosik/sso/internal/api"
	authentication "github.com/vishenosik/sso/internal/api/authentication/grpc"
	"github.com/vishenosik/sso/internal/services"
	"github.com/vishenosik/sso/internal/store/sql/sqlite"

	// internal
	embed "github.com/vishenosik/sso"
)

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type App struct {
	Server
}

func NewApp() (*App, error) {

	// Cache init

	cacheStore, err := cache.NewRedisCache()
	if err != nil {
		cacheStore = cache.NewNoopCache()
	}

	// Stores init

	sqlStore, err := sql.NewSqliteStore(
		sql.WithMigration(
			embed.Migrations,
			path.Join(embed.MigrationsPath, "sqlite"),
		),
	)
	if err != nil {
		return nil, err
	}

	db, err := sqlStore.Open(context.TODO())
	if err != nil {
		return nil, err
	}

	usersStore := sqlite.NewUsersStore(db)

	appsStore := sqlite.NewAppsStore(db)

	// Usecases init
	authService, err := services.NewAuthenticationService(usersStore, usersStore, appsStore)
	if err != nil {
		return nil, err
	}

	// Apis init

	// Services init
	handler, err := _http.NewHttpServer(newHandler(
		api.NewAuthenticationHttpApi(authService),
	))
	if err != nil {
		return nil, err
	}

	grpcServer := authentication.NewAuthenticationServer(authService)

	rpc, err := grpc.NewGrpcServer([]grpc.GrpcService{
		grpcServer,
	})
	if err != nil {
		return nil, err
	}

	app, err := gocherry.NewApp()
	if err != nil {
		return nil, err
	}

	app.AddServices(handler, cacheStore, sqlStore, rpc)

	// // Data schemas init
	// cachedStore := combined.NewCachedStore(store, cache)

	// dgraphStore, err := loadDgraph(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// // Services init
	// authenticationService := authenticationService.NewService(
	// 	log,
	// 	authenticationService.Config{
	// 		TokenTTL: conf.AuthenticationService.TokenTTL,
	// 	},
	// 	dgraphStore,
	// 	dgraphStore,
	// 	cachedStore,
	// )

	// grpcServer := grpcApp.NewGrpcApp(
	// 	log,
	// 	grpcApp.Config{
	// 		Server: config.Server{
	// 			Port: conf.GrpcConfig.Port,
	// 		},
	// 	},
	// 	authenticationService,
	// )

	// restServer := restApp.NewRestApp(
	// 	ctx,
	// 	restApp.Config{
	// 		Server: config.Server{
	// 			Port: conf.RestConfig.Port,
	// 		},
	// 	},
	// 	authenticationService,
	// )

	return &App{
		Server: app,
	}, nil
}

type Service interface {
	Routers(r chi.Router)
}

func newHandler(services ...Service) http.Handler {

	router := chi.NewRouter()
	router.Use(
		_http.SetHeaders(),
		_http.RequestLogger(),
	)

	router.Get("/swagger/*", httpSwagger.Handler())

	router.Route("/api", func(r chi.Router) {
		for i := range services {
			services[i].Routers(r)
		}
	})
	return router
}
