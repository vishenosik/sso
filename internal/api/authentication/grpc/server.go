package authentication

import (
	"context"
	"log/slog"

	authentication_v1 "github.com/blacksmith-vish/sso/gen/v1/authentication"

	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
)

type Authentication interface {
	Login(
		ctx context.Context,
		request models.LoginRequest,
	) (response models.LoginResponse, err error)

	RegisterNewUser(
		ctx context.Context,
		request models.RegisterRequest,
	) (response models.RegisterResponse, err error)

	IsAdmin(
		ctx context.Context,
		request models.IsAdminRequest,
	) (response models.IsAdminResponse, err error)
}

type authenticationAPI struct {
	authentication_v1.UnimplementedAuthenticationServer
	log  *slog.Logger
	auth Authentication
}

type server = authenticationAPI

func NewAuthenticationServer(
	log *slog.Logger,
	auth Authentication,
) *authenticationAPI {

	return &authenticationAPI{
		log:  log,
		auth: auth,
	}

}
