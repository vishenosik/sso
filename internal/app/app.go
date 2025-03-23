package app

import (
	"context"
	"fmt"
	"log/slog"

	grpcApp "github.com/vishenosik/sso/internal/app/grpc"
	restApp "github.com/vishenosik/sso/internal/app/rest"
	authenticationService "github.com/vishenosik/sso/internal/services/authentication"
	"github.com/vishenosik/sso/internal/store/combined"

	appctx "github.com/vishenosik/sso/internal/app/context"
	"github.com/vishenosik/sso/pkg/helpers/config"
)

type App struct {
	log     *slog.Logger
	servers []Server
}

type Server interface {
	MustRun()
	Stop(ctx context.Context)
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
	appContext := appctx.AppCtx(ctx)

	log := appContext.Logger
	conf := appContext.Config

	// Cache init
	cache := loadCache(ctx)

	// Stores init
	store, err := loadSqlStore(ctx)
	if err != nil {
		return nil, err
	}

	// Data schemas init
	cachedStore := combined.NewCachedStore(store, cache)

	dgraphStore, err := loadDgraph(ctx)
	if err != nil {
		return nil, err
	}

	// Services init
	authenticationService := authenticationService.NewService(
		log,
		authenticationService.Config{
			TokenTTL: conf.AuthenticationService.TokenTTL,
		},
		dgraphStore,
		dgraphStore,
		cachedStore,
	)

	grpcServer := grpcApp.NewGrpcApp(
		log,
		grpcApp.Config{
			Server: config.Server{
				Port: conf.GrpcConfig.Port,
			},
		},
		authenticationService,
	)

	restServer := restApp.NewRestApp(
		ctx,
		restApp.Config{
			Server: config.Server{
				Port: conf.RestConfig.Port,
			},
		},
		authenticationService,
	)

	return newApp(log, grpcServer, restServer), nil
}

func newApp(
	logger *slog.Logger,
	apps ...Server,
) *App {
	return &App{
		log:     logger,
		servers: apps,
	}
}

func (app *App) MustRun() {

	app.log.Info("start app")

	for _, server := range app.servers {
		go server.MustRun()
	}
}

func (app *App) Stop(ctx context.Context) {

	const msg = "app stopping"

	signal, ok := appctx.SignalCtx(ctx)
	if !ok {
		app.log.Info(msg, slog.String("signal", signal.Signal.String()))
	} else {
		app.log.Info(msg)
	}

	for _, server := range app.servers {
		server.Stop(ctx)
	}

	app.log.Info("app stopped")
}
