package restApp

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	authentication "github.com/blacksmith-vish/sso/internal/api/authentication/rest"
	_ "github.com/blacksmith-vish/sso/internal/gen/swagger"
	"github.com/blacksmith-vish/sso/internal/lib/config"
	middleW "github.com/blacksmith-vish/sso/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type App struct {
	log    *slog.Logger
	server *http.Server
	port   uint16
}

type Service interface {
	InitRouters(router *chi.Mux)
}

func NewRestApp(
	log *slog.Logger,
	conf config.RESTConfig,
	authenticationService authentication.Authentication,
) *App {

	authentication := authentication.NewAuthenticationServer(log, authenticationService)

	router := chi.NewRouter()
	router.Use(
		middleW.RequestLogger(log),
	)

	router.Get("/swagger/*", httpSwagger.Handler())

	setRouters(
		router,
		authentication,
	)

	return &App{
		log: log,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", conf.Port),
			Handler: router,
		},
		port: conf.Port,
	}
}

func setRouters(router *chi.Mux, services ...Service) {
	for i := range services {
		services[i].InitRouters(router)
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {

	const op = "restApp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Any("port", a.port),
	)

	log.Info("REST server is running", slog.String("addr", a.server.Addr))

	if err := a.server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return errors.Wrap(err, op)
		}
	}

	return nil
}

func (a *App) Stop(ctx context.Context) {

	const op = "restApp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping Rest server", slog.Any("port", a.port))

	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
}
