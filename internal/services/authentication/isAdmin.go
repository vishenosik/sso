package authentication

import (
	"context"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var (
	ErrInvalidUserID = errors.New("user_id is not valid")
)

func (a *Authentication) IsAdmin(
	ctx context.Context,
	userID string,
) (bool, error) {

	const op = "Authentication.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.String("userID", userID),
	)

	log.Info("checking if user is admin")

	if err := validator.New().Var(userID, "required,uuid4"); err != nil {
		return false, ErrInvalidUserID
	}

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		log.Error("error occured", slog.String("", err.Error()))
		return false, errors.Wrap(err, op)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
