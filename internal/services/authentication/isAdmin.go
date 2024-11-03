package authentication

import (
	"context"
	"errors"
	"log/slog"

	"github.com/blacksmith-vish/sso/internal/lib/logger"
	"github.com/blacksmith-vish/sso/internal/lib/operation"
	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
	store_models "github.com/blacksmith-vish/sso/internal/store/models"
	"github.com/go-playground/validator/v10"
)

// IsAdmin checks if user is admin
//
//	@param ctx
//	@param userID - uuid v4 ID
//
// Returned errors:
//
//	ErrUserNotFound - user not found
//	ErrUsersStore - other users store errors
//	ErrInvalidUserID - userID is invalid (basically not uuid4)
func (a *Authentication) IsAdmin(
	ctx context.Context,
	userID string,
) (bool, error) {

	var (
		op  = op("IsAdmin")
		ret = operation.ReturnFailWithError(false, op)
		log = a.log.With(
			slog.String("op", op),
			slog.String("userID", userID),
		)
	)

	log.Info("checking if user is admin")

	if err := validator.New().Var(userID, "required,uuid4"); err != nil {
		return ret(models.ErrInvalidUserID)
	}

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {

		log.Error("error occured", logger.Error(err))

		if errors.Is(err, store_models.ErrNotFound) {
			return ret(models.ErrUserNotFound)
		}

		return ret(models.ErrUsersStore)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
