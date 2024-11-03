package authentication

import (
	"context"
	"log/slog"

	authentication_v1 "github.com/blacksmith-vish/sso/gen/v1/authentication"
	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (srv *server) IsAdmin(
	ctx context.Context,
	request *authentication_v1.IsAdminRequest,
) (*authentication_v1.IsAdminResponse, error) {

	const message = "login failed"

	log := srv.log.With(
		slog.String("op", authentication_v1.Authentication_IsAdmin_FullMethodName),
		slog.String("userID", request.GetUserId()),
	)

	isAdmin, err := srv.auth.IsAdmin(
		ctx,
		request.GetUserId(),
	)

	if err != nil {

		log.Error(message, slog.String("err", err.Error()))

		if errors.Is(err, models.ErrInvalidUserID) {
			return nil, status.Error(codes.InvalidArgument, message)
		}
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, status.Error(codes.AlreadyExists, message)
		}
		return nil, status.Error(codes.Internal, message)
	}

	response := &authentication_v1.IsAdminResponse{
		IsAdmin: isAdmin,
	}

	return response, nil
}
