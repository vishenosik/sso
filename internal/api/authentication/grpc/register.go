package authentication

import (
	"context"
	"log/slog"

	authentication_v1 "github.com/blacksmith-vish/sso/internal/gen/grpc/v1/authentication"
	"github.com/blacksmith-vish/sso/internal/services/authentication/models"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *server) Register(
	ctx context.Context,
	request *authentication_v1.RegisterRequest,
) (*authentication_v1.RegisterResponse, error) {

	log := srv.log.With(
		slog.String("op", authentication_v1.Authentication_Register_FullMethodName),
	)

	serviceRequest := models.RegisterRequest{
		Nickname: "me",
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}

	if err := validator.New().Struct(serviceRequest); err != nil {
		log.Error("validation failed", "err", err.Error())
		return nil, status.Error(codes.InvalidArgument, "login failed")
	}

	UserID, err := srv.auth.RegisterNewUser(
		ctx,
		serviceRequest,
	)
	if err != nil {
		if errors.Is(err, models.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "login failed")
		}
		return nil, status.Error(codes.Internal, "login failed")
	}

	response := &authentication_v1.RegisterResponse{
		UserId: UserID,
	}

	return response, nil
}
