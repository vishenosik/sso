package authentication

import (
	"context"

	"github.com/blacksmith-vish/sso/internal/lib/logger/attrs"
	"github.com/blacksmith-vish/sso/internal/lib/operation"
	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
	store_models "github.com/blacksmith-vish/sso/internal/store/models"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

// RegisterNewUser registers new user
//
//	@param ctx
//	@param request - user data passed from registration
//
// Returned errors:
//
//	ErrInvalidRequest - one or more `@request` fields are not valid
//	ErrPasswordTooLong - password is longer than 72 bytes
//	ErrGenerateHash - failed to generate pass hash
//	ErrUserExists - user already exists
//	ErrUsersStore - other users store errors
func (auth *Authentication) RegisterNewUser(
	ctx context.Context,
	request models.RegisterRequest,
) (string, error) {

	fail, attr := operation.FailResultWithAttr("", op("Login"))
	log := auth.log.With(attr)

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
