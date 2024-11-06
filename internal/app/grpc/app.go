package grpcApp

import (
	"fmt"
	"log/slog"
	"net"

	authentication "github.com/blacksmith-vish/sso/internal/api/authentication/grpc"
	"github.com/blacksmith-vish/sso/internal/lib/config"

	authentication_v1 "github.com/blacksmith-vish/sso/sdk/api/grpc/v1/authentication"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       uint16
}

func NewGrpcApp(
	log *slog.Logger,
	conf config.GRPCConfig,
	authService authentication.Authentication,
) *App {

	Log := log.WithGroup(
		"gRPC",
	)

	gRPCServer := grpc.NewServer()

	authentication_v1.RegisterAuthenticationServer(
		gRPCServer,
		authentication.NewAuthenticationServer(
			Log,
			authService,
		),
	)

	return &App{
		log:        Log,
		gRPCServer: gRPCServer,
		port:       conf.Port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {

	const op = "grpcApp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Any("port", a.port),
	)

	log.Info("starting gRPC server")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return errors.Wrap(err, op)
	}

	log.Info("gRPC server is running", slog.String("addr", listener.Addr().String()))

	if err := a.gRPCServer.Serve(listener); err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (a *App) Stop() {

	const op = "grpcApp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Any("port", a.port))

	a.gRPCServer.GracefulStop()

}
