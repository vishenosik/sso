package authentication

import (
	// std
	"context"
	"log/slog"

	// pkg

	// internal
	authentication_v1 "github.com/vishenosik/sso/internal/gen/grpc/v1/authentication"
	"github.com/vishenosik/sso/internal/services/authentication/models"
	"github.com/vishenosik/sso/pkg/helpers/operation"
	"github.com/vishenosik/sso/pkg/logger/attrs"
)

// compileIsAdmin compiles the isAdmin function, which checks if a user is an admin.
// It takes a logger and a server instance as parameters and returns a function that can be used to check if a user is an admin.
//
// Parameters:
//
//	logger: A logger instance to log any errors that occur during the admin check.
//	srv: A server instance that has access to the authentication service.
//
// Returns:
//
//	isAdminFunc: A function that takes a context and an IsAdminRequest as parameters and returns an IsAdminResponse and an error.
func compileIsAdmin(
	logger *slog.Logger,
	srv *server,
) isAdminFunc {

	const message = "admin check failed"
	fail := operation.FailWrapErrorStatus((*authentication_v1.IsAdminResponse)(nil), message)

	return func(ctx context.Context, request *authentication_v1.IsAdminRequest) (*authentication_v1.IsAdminResponse, error) {

		log := logger.With(
			attrs.Operation(authentication_v1.Authentication_IsAdmin_FullMethodName),
			attrs.UserID(request.GetUserId()),
		)

		isAdmin, err := srv.auth.IsAdmin(ctx, request.GetUserId())
		if err != nil {
			log.Error(message, attrs.Error(err))
			return fail(models.ServiceErrorsToGrpcCodes.Get(err))
		}

		return &authentication_v1.IsAdminResponse{
			IsAdmin: isAdmin,
		}, nil
	}

}
