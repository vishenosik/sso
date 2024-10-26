package authentication

import (
	"context"
	"log/slog"

	authentication_v1 "github.com/blacksmith-vish/sso/gen/v1/authentication"
	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
	auth_store "github.com/blacksmith-vish/sso/internal/store/sql/authentication"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *server) IsAdmin(
	ctx context.Context,
	request *authentication_v1.IsAdminRequest,
) (*authentication_v1.IsAdminResponse, error) {

	log := srv.log.With(
		slog.String("op", authentication_v1.Authentication_IsAdmin_FullMethodName),
		slog.String("userID", request.GetUserId()),
	)

	serviceRequest := models.IsAdminRequest{
		UserID: request.GetUserId(),
	}

	if err := validator.New().Struct(serviceRequest); err != nil {
		log.Error("validation failed", "err", err.Error())
		return nil, status.Error(codes.InvalidArgument, "login failed")
	}

	serviceResponse, err := srv.auth.IsAdmin(
		ctx,
		serviceRequest,
	)

	if err != nil {
		if errors.Is(err, auth_store.ErrUserNotFound) {
			return nil, status.Error(codes.AlreadyExists, "login failed")
		}
		return nil, status.Error(codes.Internal, "login failed")
	}

	response := &authentication_v1.IsAdminResponse{
		IsAdmin: serviceResponse.IsAdmin,
	}

	return response, nil
}
