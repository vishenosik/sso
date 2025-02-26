package authentication

import (
	"context"
	"errors"
	"log/slog"

	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
	store_models "github.com/blacksmith-vish/sso/internal/store/models"
	"github.com/blacksmith-vish/sso/pkg/helpers/operation"
	"github.com/blacksmith-vish/sso/pkg/logger/attrs"
	"github.com/blacksmith-vish/sso/pkg/validator"
)

func compileIsAdmin(
	logger *slog.Logger,
) isAdminFunc {

	method := op("IsAdmin")
	fail := operation.FailWrapError(false, method)

	return func(ctx context.Context, auth *Authentication, userID string) (bool, error) {

		log := logger.With(
			attrs.Operation(method),
			attrs.UserID(userID),
		)

		log.Info("checking if user is admin")

		if err := validator.UUID4(userID); err != nil {
			log.Error("userID validation failed", attrs.Error(err))
			return fail(models.ErrUserInvalidID)
		}

		isAdmin, err := auth.userProvider.IsAdmin(ctx, userID)
		if err != nil {

			log.Error("error occured", attrs.Error(err))

			if errors.Is(err, store_models.ErrNotFound) {
				return fail(models.ErrUserNotFound)
			}

			return fail(models.ErrUsersStore)
		}

		log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

		return isAdmin, nil
	}

}
