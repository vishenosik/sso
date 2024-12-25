package authentication

import (
	"context"
	"log/slog"

	"github.com/blacksmith-vish/sso/internal/lib/helpers/operation"
	"github.com/blacksmith-vish/sso/internal/lib/logger/attrs"
	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
	store_models "github.com/blacksmith-vish/sso/internal/store/models"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

func compileRegisterNewUser(
	logger *slog.Logger,
) registerNewUserFunc {

	method := op("RegisterNewUser")
	fail := operation.FailResult("", method)
	log := logger.With(
		attrs.Operation(method),
	)

	return func(ctx context.Context, auth *Authentication, request models.RegisterRequest) (string, error) {

		if err := validator.New().Struct(request); err != nil {
			log.Error("failed to validate request body", attrs.Error(err))
			return fail(models.ErrInvalidRequest)
		}

		log.Info("registering user")

		passHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {

			log.Error("failed to generate pass hash", attrs.Error(err))

			if errors.Is(err, bcrypt.ErrPasswordTooLong) {
				return fail(models.ErrPasswordTooLong)
			}
			return fail(models.ErrGenerateHash)
		}

		log.Debug("generated password hash")

		userID := uuid.New().String()

		if err := auth.userSaver.SaveUser(ctx, userID, request.Nickname, request.Email, passHash); err != nil {

			log.Error("failed to save user", attrs.Error(err))

			if errors.Is(err, store_models.ErrAlreadyExists) {
				return fail(models.ErrUserExists)
			}
			return fail(models.ErrUsersStore)
		}

		log.Info("user registered successfuly")

		return userID, nil
	}

}
