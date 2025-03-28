package authentication

import (
	// std
	"context"
	"log/slog"

	// internal
	authentication_v1 "github.com/vishenosik/sso/internal/gen/grpc/v1/authentication"
	"github.com/vishenosik/sso/internal/services/authentication/models"
	"github.com/vishenosik/sso/pkg/helpers/operation"
	"github.com/vishenosik/sso/pkg/logger/attrs"
)

// compileRegisterNewUserFunc creates and returns a function for registering a new user.
// It sets up error handling, logging, and maps errors to appropriate gRPC status codes.
//
// Parameters:
//
//	logger: A pointer to a slog.Logger for logging operations.
//	srv: A pointer to the server struct containing authentication service.
//
// Returns:
//
//	registerNewUserFunc: A function that handles the registration of a new user.
//	 This returned function takes a context and a RegisterRequest, and returns
//	 a RegisterResponse and an error.
func compileRegisterNewUserFunc(
	logger *slog.Logger,
	srv *server,
) registerNewUserFunc {

	const message = "user registration failed"
	fail := operation.FailWrapErrorStatus((*authentication_v1.RegisterResponse)(nil), message)

	log := logger.With(
		attrs.Operation(authentication_v1.Authentication_Register_FullMethodName),
	)

	return func(ctx context.Context, request *authentication_v1.RegisterRequest) (*authentication_v1.RegisterResponse, error) {

		serviceRequest := models.RegisterRequest{
			Nickname: "me",
			Email:    request.GetEmail(),
			Password: request.GetPassword(),
		}

		userID, err := srv.auth.RegisterNewUser(ctx, serviceRequest)
		if err != nil {
			log.Error(message, attrs.Error(err))
			return fail(models.ServiceErrorsToGrpcCodes.Get(err))
		}

		return &authentication_v1.RegisterResponse{
			UserId: userID,
		}, nil
	}
}
