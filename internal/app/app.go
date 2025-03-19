package app

import (
	"context"
	"fmt"
	"log/slog"

	grpcApp "github.com/blacksmith-vish/sso/internal/app/grpc"
	restApp "github.com/blacksmith-vish/sso/internal/app/rest"
	authenticationService "github.com/blacksmith-vish/sso/internal/services/authentication"
	"github.com/blacksmith-vish/sso/internal/store/combined"

	appctx "github.com/blacksmith-vish/sso/internal/app/context"
	"github.com/blacksmith-vish/sso/pkg/helpers/config"
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

	// logger setup
	// TODO: implement env logic
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

	_, err = loadDgraph(ctx)
	if err != nil {
		return nil, err
	}

	// Services init
	authenticationService := authenticationService.NewService(
		log,
		authenticationService.Config{
			TokenTTL: conf.AuthenticationService.TokenTTL,
		},
		store,
		store,
		cachedStore,
	)

	return newApp(
		log,
		grpcApp.NewGrpcApp(
			log,
			grpcApp.Config{
				Server: config.Server{
					Port: conf.GrpcConfig.Port,
				},
			},
			authenticationService,
		),
		restApp.NewRestApp(
			ctx,
			restApp.Config{
				Server: config.Server{
					Port: conf.RestConfig.Port,
				},
			},
			authenticationService,
		),
	), nil
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
