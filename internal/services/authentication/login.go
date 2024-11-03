package authentication

import (
	"context"
	"log/slog"

	"errors"

	"github.com/blacksmith-vish/sso/internal/lib/jwt"
	"github.com/blacksmith-vish/sso/internal/lib/logger"
	"github.com/blacksmith-vish/sso/internal/lib/operation"
	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
	store_models "github.com/blacksmith-vish/sso/internal/store/models"
	"github.com/go-playground/validator/v10"

	"golang.org/x/crypto/bcrypt"
)

// Login checks if user's credentials exists and appID is valid
//
//	@param ctx
//	@param request - user data passed from login
//	@param appID - uuid v4 ID
//
// Returned errors:
//
//	ErrInvalidRequest - one or more `@request` fields are not valid
//	ErrInvalidAppID - appID is invalid (basically not uuid4)
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

	const noToken = ""

	var (
		op    = op("Login")
		ret   = operation.ReturnFailWithError(noToken, op)
		valid = validator.New()
		log   = auth.log.With(
			slog.String("op", op),
			slog.String("app_id", appID),
		)
	)

	if err := valid.Var(appID, "required,uuid4"); err != nil {
		log.Error("appID validation failed", logger.Error(err))
		return ret(models.ErrInvalidAppID)
	}

	if err := valid.Struct(request); err != nil {
		log.Error("failed to validate request body", logger.Error(err))
		return ret(models.ErrInvalidRequest)
	}

	log.Debug("attempting to get app")

	app, err := auth.appProvider.App(ctx, appID)
	if err != nil {

		log.Error("failed to get app", logger.Error(err))

		if errors.Is(err, store_models.ErrNotFound) {
			return ret(models.ErrAppNotFound)
		}

		return ret(models.ErrAppsStore)
	}

	log.Info("attempting to login user")

	user, err := auth.userProvider.UserByEmail(ctx, request.Email)
	if err != nil {
		log.Error("failed to get user", logger.Error(err))

		if errors.Is(err, store_models.ErrNotFound) {
			return ret(models.ErrUserNotFound)
		}
		return ret(models.ErrUsersStore)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(request.Password)); err != nil {
		log.Error(
			"invalid password",
			slog.String("err", err.Error()),
			slog.String("user_id", user.ID),
		)
		return ret(models.ErrInvalidCredentials)
	}

	token := jwt.NewToken(user, app, auth.tokenTTL)

	log.Info("user logged in succesfully", slog.String("user_id", user.ID))

	return token, nil
}
