package authentication

import (
	"context"
	"log/slog"

	authentication_v1 "github.com/blacksmith-vish/sso/gen/v1/authentication"
	"github.com/blacksmith-vish/sso/internal/services/authentication"
	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
	"github.com/pkg/errors"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *server) Login(
	ctx context.Context,
	request *authentication_v1.LoginRequest,
) (*authentication_v1.LoginResponse, error) {

	log := srv.log.With(
		slog.String("op", authentication_v1.Authentication_Login_FullMethodName),
	)

	serviceRequest := models.LoginRequest{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}

	if err := validator.New().Struct(serviceRequest); err != nil {
		log.Error("validation failed", "err", err.Error())
		return nil, status.Error(codes.InvalidArgument, "login failed")
	}

	token, err := srv.auth.Login(
		ctx,
		serviceRequest,
		request.GetAppId(),
	)
	if err != nil {

		if errors.Is(err, authentication.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "login failed")
		}

		return nil, status.Error(codes.Internal, "login failed")
	}

	response := &authentication_v1.LoginResponse{
		Token: token,
	}

	return response, nil
}
