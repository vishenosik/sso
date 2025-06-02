package entities

import (
	"github.com/pkg/errors"
	"github.com/vishenosik/sso/pkg/helpers/errorHelper"
	"google.golang.org/grpc/codes"
)

var (
	// no content added
	ErrContentNotAdded = errors.New("content is not added")
)

// validation errors
var (
	// invalid ID
	ErrID = errors.New("invalid ID")
	// invalid email
	ErrEmail = errors.New("invalid email")
	// provided string is not URL
	ErrURL = errors.New("provided string is not URL")
	// time interval can't be less than time.Minute
	ErrInterval = errors.New("time interval can't be less than time.Minute")
	// string must consist of only ascii characters
	ErrAscii = errors.New("string must consist of only ascii characters")
	// must be in (0,600) interval
	ErrCode = errors.New("must be in (0,600) interval")
	// is required
	ErrRequired = errors.New("is required")
)

var (
	// COMMON

	// request fields are invalid
	ErrInvalidRequest = errors.New("failed to validate request body")

	// APP

	// app not found
	ErrAppNotFound = errors.New("app not found")
	// invalid app_id
	ErrAppInvalidID = errors.New("invalid app_id")
	// apps store unexpected error
	ErrAppsStore = errors.New("apps store unexpected error")

	// USER

	// user exists already
	ErrUserExists = errors.New("user exists already")
	// user not found
	ErrUserNotFound = errors.New("user not found")
	// invalid credentials
	ErrInvalidCredentials = errors.New("invalid credentials")
	// invalid user_id
	ErrUserInvalidID = errors.New("invalid user_id")
	// failed to generate pass hash
	ErrGenerateHash = errors.New("failed to generate pass hash")
	// password length exceeds 72 bytes
	ErrPasswordTooLong = errors.New("password length exceeds 72 bytes")
	// users store unexpected error
	ErrUsersStore = errors.New("users store unexpected error")
)

var (
	// not found
	ErrNotFound = errors.New("not found")
	// exists already
	ErrAlreadyExists = errors.New("exists already")
)

var ServiceErrorsToGrpcCodes *errorHelper.ErrorsMap[codes.Code]

func init() {
	ServiceErrorsToGrpcCodes = errorHelper.NewErrorsMap(
		map[error]codes.Code{
			ErrInvalidRequest:     codes.InvalidArgument,
			ErrAppNotFound:        codes.NotFound,
			ErrAppInvalidID:       codes.InvalidArgument,
			ErrAppsStore:          codes.Internal,
			ErrUserExists:         codes.AlreadyExists,
			ErrUsersStore:         codes.Internal,
			ErrUserInvalidID:      codes.InvalidArgument,
			ErrUserNotFound:       codes.NotFound,
			ErrInvalidCredentials: codes.Unauthenticated,
			ErrGenerateHash:       codes.Internal,
			ErrPasswordTooLong:    codes.InvalidArgument,
		},
		codes.Internal,
	)
}
