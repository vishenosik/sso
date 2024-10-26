package authentication_v1

import (
	"context"
	"log/slog"

	"github.com/blacksmith-vish/sso/pkg/api"
	"github.com/go-chi/chi/v5"
)

const (
	routeUser string = "user"
)

type Authentication interface {
	Register(
		ctx context.Context,
		request RegisterRequest,
	) (RegisterResponse, error)
}

type authenticationServer struct {
	log  *slog.Logger
	auth Authentication
}

type server = *authenticationServer

func NewAuthenticationServer(
	log *slog.Logger,
	auth Authentication,
) *authenticationServer {
	return &authenticationServer{
		log:  log,
		auth: auth,
	}
}

func (srv server) RegisterAuthenticationServer(mux *chi.Mux) {

	router := chi.NewRouter()

	router.Post(pathRegister, srv.register())

	// Mounting the new sub routers on the main router
	mux.Mount(
		api.ApiV1(routeUser),
		router,
	)
}
