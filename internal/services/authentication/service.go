// Authentication package provides business methods to identify/save/check users
//
//go:generate go run github.com/vektra/mockery/v2@v2.45.0 --all --case=camel
package authentication

import (
	"context"
	"log/slog"
	"time"

	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
	store_models "github.com/blacksmith-vish/sso/internal/store/models"
	"github.com/blacksmith-vish/sso/pkg/helpers/operation"
)

const (
	userSaver_SaveUser       = "SaveUser"
	userProvider_UserByEmail = "UserByEmail"
	userProvider_IsAdmin     = "IsAdmin"
	appProvider_App          = "App"
)

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		id string,
		nickname string,
		email string,
		passwordHash []byte,
	) (err error)
}

type UserProvider interface {
	UserByEmail(
		ctx context.Context,
		email string,
	) (user store_models.User, err error)
	IsAdmin(
		ctx context.Context,
		userID string,
	) (isAdmin bool, err error)
}

type AppProvider interface {
	App(
		ctx context.Context,
		appID string,
	) (app store_models.App, err error)
}

type (
	isAdminFunc = func(
		ctx context.Context,
		auth *Authentication,
		userID string,
	) (bool, error)

	registerNewUserFunc = func(
		ctx context.Context,
		auth *Authentication,
		request models.RegisterRequest,
	) (string, error)

	loginFunc = func(
		ctx context.Context,
		auth *Authentication,
		request models.LoginRequest,
		appID string,
	) (string, error)
)

var (
	isAdmin         isAdminFunc
	registerNewUser registerNewUserFunc
	login           loginFunc
)

type Authentication struct {
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type Config struct {
	TokenTTL time.Duration
}

// New returns a new instance of Auth
func NewService(
	logger *slog.Logger,
	conf Config,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
) *Authentication {

	isAdmin = compileIsAdmin(logger)
	login = compileLogin(logger)
	registerNewUser = compileRegisterNewUser(logger)

	return &Authentication{
		tokenTTL:     conf.TokenTTL,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
	}
}

// IsAdmin checks if user is admin
//
//	param ctx
//	param userID - uuid v4 ID
//
// Returned errors:
//
//	ErrUserNotFound - user not found
//	ErrUsersStore - other users store errors
//	ErrUserInvalidID - userID is invalid (basically not uuid4)
func (auth *Authentication) IsAdmin(
	ctx context.Context,
	userID string,
) (bool, error) {
	return isAdmin(ctx, auth, userID)
}

// RegisterNewUser registers new user
//
//	param ctx
//	param request - user data passed from registration
//
// Returned errors:
//
//	ErrInvalidRequest - one or more `@request` fields are not valid
//	ErrPasswordTooLong - password is longer than 72 bytes
//	ErrGenerateHash - failed to generate pass hash
//	ErrUserExists - user already exists
//	ErrUsersStore - other users store errors
func (auth *Authentication) RegisterNewUser(
	ctx context.Context,
	request models.RegisterRequest,
) (string, error) {
	return registerNewUser(ctx, auth, request)
}

// Login checks if user's credentials exists and appID is valid
//
//	param ctx
//	param request - user data passed from login
//	param appID - uuid v4 ID
//
// Returned errors:
//
//	ErrInvalidRequest - one or more `@request` fields are not valid
//	ErrAppInvalidID - appID is invalid (basically not uuid4)
//	ErrAppNotFound - app not found
//	ErrAppsStore - other apps store errors
//	ErrUserNotFound - user not found
//	ErrUsersStore - other users store errors
//	ErrInvalidCredentials - invalid password passed
func (auth *Authentication) Login(
	ctx context.Context,
	request models.LoginRequest,
	appID string,
) (string, error) {
	return login(ctx, auth, request, appID)
}

func op(method string) string {
	return operation.ServicesOperation("authentication", method)
}
