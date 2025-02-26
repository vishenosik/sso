package authentication

import (
	"context"
	"log/slog"

	"errors"

	"github.com/blacksmith-vish/sso/internal/lib/jwt"
	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
	store_models "github.com/blacksmith-vish/sso/internal/store/models"
	"github.com/blacksmith-vish/sso/pkg/helpers/operation"
	"github.com/blacksmith-vish/sso/pkg/logger/attrs"
	"github.com/blacksmith-vish/sso/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

func compileLogin(
	logger *slog.Logger,
) loginFunc {

	method := op("Login")
	fail := operation.FailWrapError("", method)

	return func(ctx context.Context, auth *Authentication, request models.LoginRequest, appID string) (string, error) {

		log := logger.With(
			attrs.Operation(method),
			attrs.AppID(appID),
		)

		if err := validator.UUID4(appID); err != nil {
			log.Error("appID validation failed", attrs.Error(err))
			return fail(models.ErrAppInvalidID)
		}

		if err := validator.Struct(request); err != nil {
			log.Error("failed to validate request body", attrs.Error(err))
			return fail(models.ErrInvalidRequest)
		}

		log.Debug("attempting to get app")

		app, err := auth.appProvider.App(ctx, appID)
		if err != nil {

			log.Error("failed to get app", attrs.Error(err))

			if errors.Is(err, store_models.ErrNotFound) {
				return fail(models.ErrAppNotFound)
			}

			return fail(models.ErrAppsStore)
		}

		log.Info("attempting to login user")

		user, err := auth.userProvider.UserByEmail(ctx, request.Email)
		if err != nil {
			log.Error("failed to get user", attrs.Error(err))

			if errors.Is(err, store_models.ErrNotFound) {
				return fail(models.ErrUserNotFound)
			}
			return fail(models.ErrUsersStore)
		}

		userIDAttr := attrs.UserID(user.ID)

		if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(request.Password)); err != nil {
			log.Error("invalid password", attrs.Error(err), userIDAttr)
			return fail(models.ErrInvalidCredentials)
		}

		token := jwt.NewToken(user, app, auth.tokenTTL)

		log.Info("user logged in succesfully", userIDAttr)

		return token, nil
	}

}
