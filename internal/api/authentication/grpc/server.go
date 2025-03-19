//go:generate go run github.com/vektra/mockery/v2v2.45.0:-all:-case=camel
package authentication

import (
	// std
	"context"
	"log/slog"

	// internal
	authentication_v1 "github.com/vishenosik/sso/internal/gen/grpc/v1/authentication"
	"github.com/vishenosik/sso/internal/services/authentication/models"
)

type Authentication interface {
	Login(
		ctx context.Context,
		request models.LoginRequest,
		appID string,
	) (token string, err error)

	RegisterNewUser(
		ctx context.Context,
		request models.RegisterRequest,
	) (userID string, err error)

	IsAdmin(
		ctx context.Context,
		userID string,
	) (isAdmin bool, err error)
}

type isAdminFunc = func(
	ctx context.Context,
	request *authentication_v1.IsAdminRequest,
) (*authentication_v1.IsAdminResponse, error)

type loginFunc = func(
	ctx context.Context,
	request *authentication_v1.LoginRequest,
) (*authentication_v1.LoginResponse, error)

type registerNewUserFunc = func(
	ctx context.Context,
	request *authentication_v1.RegisterRequest,
) (*authentication_v1.RegisterResponse, error)

var (
	isAdmin         isAdminFunc
	login           loginFunc
	registerNewUser registerNewUserFunc
)

type authenticationAPI struct {
	authentication_v1.UnimplementedAuthenticationServer
	auth Authentication
}

type server = authenticationAPI

// NewAuthenticationServer creates a new instance of the AuthenticationServer with the provided logger and authentication service.
//
// Parameters:
//
//	log: A logger instance to log any errors or information during the authentication process.
//	auth: An instance of the Authentication interface that provides the actual authentication logic.
//
// Returns:
//
//	srv: A new instance of the AuthenticationServer that can be used to handle authentication requests.
func NewAuthenticationServer(
	log *slog.Logger,
	auth Authentication,
) *authenticationAPI {

	srv := &authenticationAPI{
		auth: auth,
	}

	isAdmin = compileIsAdmin(log, srv)
	login = compileLoginFunc(log, srv)
	registerNewUser = compileRegisterNewUserFunc(log, srv)

	return srv
}

// IsAdmin Checks if the user with the given ID is an admin.
//
// Parameters:
//
//	ctx: The context of the request.
//	userID: The ID of the user to check.
//
// Returns:
//
//	isAdmin: A boolean value indicating if the user is an admin.
//	err: An error if an issue occurs during the check.
func (srv *server) IsAdmin(
	ctx context.Context,
	request *authentication_v1.IsAdminRequest,
) (*authentication_v1.IsAdminResponse, error) {
	return isAdmin(ctx, request)
}

// Login Logs the user in using the provided credentials.
//
// Parameters:
//
//	ctx: The context of the request.
//	request: The login request containing the user credentials.
//	appID: The ID of the application making the request.
//
// Returns:
//
//	token: string representing user token.
//	err: error if an issue occurs during the login process.
func (srv *server) Login(
	ctx context.Context,
	request *authentication_v1.LoginRequest,
) (*authentication_v1.LoginResponse, error) {
	return login(ctx, request)
}

// Register Registers a new user in the system.
//
// Parameters:
//
//	ctx: The context of the request.
//	request: The registration request containing the user information.
//
// Returns:
//
//	userID: A string representing the user ID.
//	err: An error if an issue occurs during the registration process.
func (srv *server) Register(
	ctx context.Context,
	request *authentication_v1.RegisterRequest,
) (*authentication_v1.RegisterResponse, error) {
	return registerNewUser(ctx, request)
}
