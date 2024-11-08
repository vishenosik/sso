package authentication

import (
	"context"
	"errors"
	"log/slog"

	"github.com/blacksmith-vish/sso/internal/lib/logger/attrs"
	"github.com/blacksmith-vish/sso/internal/lib/operation"
	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
	store_models "github.com/blacksmith-vish/sso/internal/store/models"
	"github.com/go-playground/validator/v10"
)

// IsAdmin checks if user is admin
//
//	param ctx
//	param userID - uuid v4 ID
//
// Returned errors:
//
//	ErrUserNotFound - user not found
//	ErrUsersStore - other users store errors
//	ErrInvalidUserID - userID is invalid (basically not uuid4)
func (auth *Authentication) IsAdmin(
	ctx context.Context,
	userID string,
) (bool, error) {

	fail, attr := operation.FailResultWithAttr(false, op("IsAdmin"))

	log := auth.log.With(
		attr,
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
