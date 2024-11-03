// Authentication package provides business methods to identify/save/check users
//
//go:generate go run github.com/vektra/mockery/v2@v2.45.0 --all --case=camel
package authentication

import (
	"context"
	"log/slog"
	"time"

	"github.com/blacksmith-vish/sso/internal/lib/config"
	"github.com/blacksmith-vish/sso/internal/lib/operation"
	"github.com/blacksmith-vish/sso/internal/store/models"
)

const (
	userSaver_SaveUser = "SaveUser"
)

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		id string,
		nickname string,
		email string,
		passwordHash []byte,
	) (err error)
}

const (
	userProvider_UserByEmail = "UserByEmail"
	userProvider_IsAdmin     = "IsAdmin"
)

type UserProvider interface {
	UserByEmail(
		ctx context.Context,
		email string,
	) (user models.User, err error)
	IsAdmin(
		ctx context.Context,
		userID string,
	) (isAdmin bool, err error)
}

const (
	appProvider_App = "App"
)

type AppProvider interface {
	App(
		ctx context.Context,
		appID string,
	) (app models.App, err error)
}

type Authentication struct {
	log          *slog.Logger
	tokenTTL     time.Duration
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
}

// New returns a new instance of Auth
func NewService(
	log *slog.Logger,
	conf config.AuthenticationService,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
) *Authentication {
	return &Authentication{
		log:          log,
		tokenTTL:     conf.TokenTTL,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
	}
}

func op(method string) string {
	return operation.ServicesOperation("Authentication", method)
}
