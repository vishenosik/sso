package main

import (
	// std
	"context"
	"flag"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	// internal pkg
	"github.com/vishenosik/gocherry"
	"github.com/vishenosik/gocherry/pkg/cache"
	_ctx "github.com/vishenosik/gocherry/pkg/context"
	"github.com/vishenosik/gocherry/pkg/grpc"
	_http "github.com/vishenosik/gocherry/pkg/http"
	"github.com/vishenosik/gocherry/pkg/sql"
	"github.com/vishenosik/sso-sdk/api"

	// internal
	embed "github.com/vishenosik/sso"
	"github.com/vishenosik/sso/internal/dto"
	"github.com/vishenosik/sso/internal/services"
	"github.com/vishenosik/sso/internal/store/sql/sqlite"
)

func main() {

	gocherry.ConfigFlags(
		services.AuthenticationConfigEnv{},
	)

	flag.Parse()
	ctx := context.Background()

	// App init
	application, err := NewApp()
	if err != nil {
		panic(err)
	}

	err = application.Start(ctx)
	if err != nil {
		panic(err)
	}

	// Graceful shut down
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	stopctx, cancel := context.WithTimeout(
		_ctx.WithStopCtx(ctx, <-stop),
		time.Second*5,
	)
	defer cancel()

	application.Stop(stopctx)
}

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

	authDTO := dto.NewAuthenticationDTO(authService)

	sys := services.NewSystem(100, false, 322)
	systemDTO := dto.NewSystemDTO(sys)
	// Apis init

	authApi := api.NewAuthenticationApi(authDTO)

	// Services init
	httpServer, err := _http.NewHttpServer(api.NewHttpHandler(
		authApi,
		api.NewSystemApi(systemDTO),
	))
	if err != nil {
		return nil, err
	}

	grpcServer, err := grpc.NewGrpcServer([]grpc.GrpcService{
		authApi,
	})
	if err != nil {
		return nil, err
	}

	app, err := gocherry.NewApp()
	if err != nil {
		return nil, err
	}

	app.AddServices(httpServer, cacheStore, sqlStore, grpcServer)

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
