// AuthenticationService package provides business methods to identify/save/check users
package services

import (
	// std
	"context"
	"log/slog"
	"time"

	// pkg

	"github.com/google/uuid"
	"github.com/vishenosik/gocherry/pkg/config"
	"github.com/vishenosik/gocherry/pkg/errors"
	"github.com/vishenosik/gocherry/pkg/logs"
	"golang.org/x/crypto/bcrypt"

	// internal pkg
	"github.com/vishenosik/gocherry/pkg/validator"

	// internal
	"github.com/vishenosik/gocherry/pkg/operation"
	"github.com/vishenosik/sso/internal/entities"
	"github.com/vishenosik/sso/internal/lib/jwt"
	store_models "github.com/vishenosik/sso/internal/store/models"
)

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		user *entities.UserCreds,
	) (err error)
}

type UserProvider interface {
	//
	UserByEmail(ctx context.Context, email string) (user *entities.UserCreds, err error)
	//
	IsAdmin(ctx context.Context, userID string) (isAdmin bool, err error)
}

type AppProvider interface {
	//
	App(ctx context.Context, appID string) (app *entities.App, err error)
}

type AuthenticationService struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider

	config AuthenticationConfig
}

type AuthenticationConfig struct {
	TokenTTL time.Duration
}

type AuthenticationConfigEnv struct {
	TokenTTL time.Duration `env:"AUTH_TOKEN_TTL" default:"1h" desc:"authentication service jwt tokens TTL"`
}

type AuthenticationServiceOption func(*AuthenticationService)

func (c AuthenticationConfig) validate() error {
	return nil
}

// New returns a new instance of Auth
func NewAuthenticationService(
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	opts ...AuthenticationServiceOption,
) (*AuthenticationService, error) {

	if userSaver == nil {
		return nil, errors.New("userSaver is nil")
	}

	if userProvider == nil {
		return nil, errors.New("userProvider is nil")
	}

	if appProvider == nil {
		return nil, errors.New("appProvider is nil")
	}

	log := logs.SetupLogger().With(
		logs.Operation("NewAuthenticationService"),
	)

	var envConf AuthenticationConfigEnv
	if err := config.ReadConfig(&envConf); err != nil {
		log.Warn("failed to read env config", logs.Error(err))
	}

	config := AuthenticationConfig{
		TokenTTL: envConf.TokenTTL,
	}

	auth := &AuthenticationService{
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		config:       config,
	}

	auth.log = log

	for _, opt := range opts {
		opt(auth)
	}

	if err := auth.config.validate(); err != nil {
		return nil, errors.Wrap(err, "failed to validate authentication service config")
	}

	return auth, nil
}

func WithConfig(config AuthenticationConfig) AuthenticationServiceOption {
	return func(as *AuthenticationService) {
		as.config = config
	}
}

func WithLogger(logger *slog.Logger) AuthenticationServiceOption {
	return func(as *AuthenticationService) {
		if logger != nil {
			as.log = logger
		}
	}
}

// RegisterUser registers new user
//
//	param ctx
//	param request - user data passed from registration
//
// Returns
//   - user's ID
//   - error
//
// Returned errors:
//
//	ErrInvalidRequest - one or more `@request` fields are not valid
//	ErrPasswordTooLong - password is longer than 72 bytes
//	ErrGenerateHash - failed to generate pass hash
//	ErrUserExists - user already exists
//	ErrUsersStore - other users store errors
func (auth *AuthenticationService) RegisterUser(
	ctx context.Context,
	user *entities.UserCreds,
) (string, error) {

	method := authenticationOP("RegisterNewUser")
	fail := errors.FailWrapError("", method)

	log := auth.log.With(
		logs.Operation(method),
	)

	user.ID = uuid.New().String()

	if err := validator.Struct(user); err != nil {
		log.Error("failed to validate request body", logs.Error(err))
		return fail(entities.ErrInvalidRequest)
	}

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {

		log.Error("failed to generate pass hash", logs.Error(err))

		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return fail(entities.ErrPasswordTooLong)
		}
		return fail(entities.ErrGenerateHash)
	}

	log.Debug("generated password hash")

	user.PasswordHash = passHash

	if err := auth.userSaver.SaveUser(ctx, user); err != nil {

		log.Error("failed to save user", logs.Error(err))

		if errors.Is(err, store_models.ErrAlreadyExists) {
			return fail(entities.ErrUserExists)
		}
		return fail(entities.ErrUsersStore)
	}

	log.Info("user registered successfuly")

	return user.ID, nil

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
func (auth *AuthenticationService) IsAdmin(
	ctx context.Context,
	userID string,
) (bool, error) {

	method := authenticationOP("IsAdmin")
	fail := errors.FailWrapError(false, method)

	log := auth.log.With(
		logs.Operation(method),
		logs.UserID(userID),
	)

	log.Info("checking if user is admin")

	if err := validator.UUID4(userID); err != nil {
		log.Error("userID validation failed", logs.Error(err))
		return fail(entities.ErrUserInvalidID)
	}

	isAdmin, err := auth.userProvider.IsAdmin(ctx, userID)
	if err != nil {

		log.Error("error occured", logs.Error(err))

		if errors.Is(err, store_models.ErrNotFound) {
			return fail(entities.ErrUserNotFound)
		}

		return fail(entities.ErrUsersStore)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}

// LoginByEmail checks if user's credentials exists and appID is valid
//
//	param ctx
//	param request - user data passed from login
//	param appID - uuid v4 ID
//
// Returns
//   - jwt token
//   - error
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
func (auth *AuthenticationService) LoginByEmail(
	ctx context.Context,
	email,
	password,
	appID string,
) (string, error) {

	method := authenticationOP("LoginByEmail")
	fail := errors.FailWrapError("", method)

	log := auth.log.With(
		logs.Operation(method),
		logs.AppID(appID),
	)

	if err := validator.UUID4(appID); err != nil {
		log.Error("appID validation failed", logs.Error(err))
		return fail(entities.ErrID)
	}

	if err := validator.Email(email); err != nil {
		log.Error("email validation failed", logs.Error(err))
		return fail(entities.ErrEmail)
	}

	if password == "" {
		log.Error("no password provided")
		return fail(errors.Wrap(entities.ErrRequired, "password"))
	}

	log.Debug("attempting to get app")

	app, err := auth.appProvider.App(ctx, appID)
	if err != nil {

		log.Error("failed to get app", logs.Error(err))

		if errors.Is(err, store_models.ErrNotFound) {
			return fail(entities.ErrAppNotFound)
		}

		return fail(entities.ErrAppsStore)
	}

	log.Info("attempting to login user")

	user, err := auth.userProvider.UserByEmail(ctx, email)
	if err != nil {
		log.Error("failed to get user", logs.Error(err))

		if errors.Is(err, store_models.ErrNotFound) {
			return fail(entities.ErrUserNotFound)
		}
		return fail(entities.ErrUsersStore)
	}

	userIDAttr := logs.UserID(user.ID)

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		log.Error("invalid password", logs.Error(err), userIDAttr)
		return fail(entities.ErrInvalidCredentials)
	}

	token := jwt.NewToken(user, app, auth.config.TokenTTL)

	log.Info("user logged in succesfully", userIDAttr)

	return token, nil
}

func authenticationOP(method string) string {
	return operation.ServicesOperation("authentication", method)
}
