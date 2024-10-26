package authentication

import (
	"context"

	"log/slog"

	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
	auth_store "github.com/blacksmith-vish/sso/internal/store/sql/authentication"

	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

func (a *Authentication) RegisterNewUser(
	ctx context.Context,
	request models.RegisterRequest,
) (models.RegisterResponse, error) {

	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", request.Email), // TODO email лучше не логировать
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate pass hash", slog.String("", err.Error()))
		return models.RegisterResponse{}, errors.Wrap(err, op)
	}

	ID, err := a.userSaver.SaveUser(ctx, request.Nickname, request.Email, passHash)
	if err != nil {

		if errors.Is(err, auth_store.ErrUserExists) {
			log.Warn("user exists", slog.String("", err.Error()))
			return models.RegisterResponse{}, errors.Wrap(auth_store.ErrUserExists, op)
		}

		log.Error("failed to save user", slog.String("", err.Error()))
		return models.RegisterResponse{}, errors.Wrap(err, op)
	}

	log.Info("user registered")

	return models.RegisterResponse{
		UserID: ID,
	}, nil

}
