package models

import "github.com/pkg/errors"

var (
	// COMMON

	// request fields are invalid
	ErrInvalidRequest = errors.New("failed to validate request body")

	// APP

	// app not found
	ErrAppNotFound = errors.New("app not found")
	// invalid app_id
	ErrInvalidAppID = errors.New("invalid app_id")
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
	ErrInvalidUserID = errors.New("invalid user_id")
	// failed to generate pass hash
	ErrGenerateHash = errors.New("failed to generate pass hash")
	// password length exceeds 72 bytes
	ErrPasswordTooLong = errors.New("password length exceeds 72 bytes")
	// users store unexpected error
	ErrUsersStore = errors.New("users store unexpected error")
)
