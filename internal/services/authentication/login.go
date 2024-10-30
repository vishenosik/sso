package authentication

import (
	"context"
	"log/slog"

	"github.com/blacksmith-vish/sso/internal/lib/jwt"
	"github.com/blacksmith-vish/sso/internal/lib/logger"
	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
	store_models "github.com/blacksmith-vish/sso/internal/store/models"
	auth_store "github.com/blacksmith-vish/sso/internal/store/sql/authentication"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

// Login checks if user's credentials exists
func (a *Authentication) Login(
	ctx context.Context,
	request models.LoginRequest,
	appID string,
) (string, error) {

	const noToken = ""

	op := op("Login")

	log := a.log.With(
		slog.String("op", op),
		slog.String("app_id", appID),
	)

	app, err := a.app(ctx, log, appID)
	if err != nil {
		return noToken, err
	}

	log.Info("attempting to login user")

	user, err := a.user(ctx, log, request)
	if err != nil {
		return noToken, err
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(request.Password)); err != nil {
		log.Error("invalid password", slog.String("err", err.Error()))
		return noToken, errors.Wrap(ErrInvalidCredentials, op)
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to get token", logger.Error(err))
		return noToken, errors.Wrap(err, op)
	}

	log.Info("user logged in succesfully", slog.String("user_id", user.ID))

	return token, nil
}

func (a *Authentication) app(
	ctx context.Context,
	log *slog.Logger,
	appID string,
) (store_models.App, error) {

	var noApp store_models.App
	op := op("app")

	if err := validator.New().Var(appID, "required,uuid4"); err != nil {
		log.Error("appID validation failed", logger.Error(err))
		return noApp, ErrInvalidAppID
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {

		if errors.Is(err, auth_store.ErrAppNotFound) {
			log.Error("app not found", logger.Error(err))
			return noApp, errors.Wrap(ErrInvalidAppID, op)
		}

		log.Error("failed to get app", logger.Error(err))
		return noApp, errors.Wrap(err, op)
	}
	return app, nil
}

func (a *Authentication) user(
	ctx context.Context,
	log *slog.Logger,
	request models.LoginRequest,
) (store_models.User, error) {

	var noUser store_models.User
	op := op("app")

	if err := validator.New().Struct(request); err != nil {
		return noUser, ErrInvalidCredentials
	}

	user, err := a.userProvider.User(ctx, request.Email)
	if err != nil {

		if errors.Is(err, auth_store.ErrUserNotFound) {
			log.Error("user not found", logger.Error(err))

			return noUser, errors.Wrap(ErrInvalidCredentials, op)
		}

		log.Error("failed to get user", logger.Error(err))
		return noUser, errors.Wrap(err, op)
	}
	return user, nil
}
