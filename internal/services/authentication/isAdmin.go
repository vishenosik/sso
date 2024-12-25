package authentication

import (
	"context"
	"errors"
	"log/slog"

	"github.com/blacksmith-vish/sso/internal/lib/helpers/operation"
	"github.com/blacksmith-vish/sso/internal/lib/logger/attrs"
	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
	store_models "github.com/blacksmith-vish/sso/internal/store/models"
	"github.com/go-playground/validator/v10"
)

func compileIsAdmin(
	logger *slog.Logger,
) isAdminFunc {

	method := op("IsAdmin")
	fail := operation.FailResult(false, method)

	return func(ctx context.Context, auth *Authentication, userID string) (bool, error) {

		log := logger.With(
			attrs.Operation(method),
			slog.String("userID", userID),
		)

		log.Info("checking if user is admin")

		if err := validator.New().Var(userID, "required,uuid4"); err != nil {
			return fail(models.ErrInvalidUserID)
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
