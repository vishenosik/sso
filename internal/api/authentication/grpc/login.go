package authentication

import (
	// std
	"context"
	"log/slog"

	// pkg

	"google.golang.org/grpc/codes"

	// internal
	authentication_v1 "github.com/blacksmith-vish/sso/internal/gen/grpc/v1/authentication"
	"github.com/blacksmith-vish/sso/internal/lib/helpers/errorHelper"
	"github.com/blacksmith-vish/sso/internal/lib/helpers/operation"
	"github.com/blacksmith-vish/sso/internal/lib/logger/attrs"
	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
)

// compileLoginFunc compiles a login function that authenticates a user and returns a token.
//
// Parameters:
//
//	logger: A logger instance to log errors and operations.
//	srv: A server instance that contains the authentication service.
//
// Returns:
//
//	loginFunc: A login function that takes a context, a LoginRequest message, and an optional AppId.
//	It returns a LoginResponse message containing a token if the login is successful,
//	or an error with the appropriate status code based on the error encountered.
func compileLoginFunc(
	logger *slog.Logger,
	srv *server,
) loginFunc {

	const message = "login failed"
	fail := operation.FailWrapErrorStatus((*authentication_v1.LoginResponse)(nil), message)

	codeMap := errorHelper.NewErrorsMap(
		map[error]codes.Code{
			models.ErrInvalidRequest:     codes.InvalidArgument,
			models.ErrInvalidAppID:       codes.InvalidArgument,
			models.ErrAppNotFound:        codes.NotFound,
			models.ErrAppsStore:          codes.Internal,
			models.ErrUserNotFound:       codes.NotFound,
			models.ErrUsersStore:         codes.Internal,
			models.ErrInvalidCredentials: codes.Unauthenticated,
		},
		codes.Internal,
	)

	log := logger.With(
		attrs.Operation(authentication_v1.Authentication_Login_FullMethodName),
	)

	return func(ctx context.Context, request *authentication_v1.LoginRequest) (*authentication_v1.LoginResponse, error) {

		serviceRequest := models.LoginRequest{
			Email:    request.GetEmail(),
			Password: request.GetPassword(),
		}

		token, err := srv.auth.Login(ctx, serviceRequest, request.GetAppId())
		if err != nil {
			log.Error(message, attrs.Error(err))
			return fail(codeMap.Get(err))
		}

		return &authentication_v1.LoginResponse{
			Token: token,
		}, nil
	}
}
