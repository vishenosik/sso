package authentication

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/vishenosik/sso/internal/services/authentication/models"
	store_models "github.com/vishenosik/sso/internal/store/models"
	"github.com/vishenosik/sso/pkg/helpers/operation"
	"github.com/vishenosik/sso/pkg/logger/attrs"
	"github.com/vishenosik/sso/pkg/validator"

	"golang.org/x/crypto/bcrypt"
)

func compileRegisterNewUser(
	logger *slog.Logger,
) registerNewUserFunc {

	method := op("RegisterNewUser")
	fail := operation.FailWrapError("", method)
	log := logger.With(
		attrs.Operation(method),
	)

	return func(ctx context.Context, auth *Authentication, request models.RegisterRequest) (string, error) {

		if err := validator.Struct(request); err != nil {
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
