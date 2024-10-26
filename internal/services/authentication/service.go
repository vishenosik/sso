package authentication

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/blacksmith-vish/sso/internal/lib/config"
	"github.com/blacksmith-vish/sso/internal/store/models"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app_id")
)

//go:generate go run github.com/vektra/mockery/v2@v2.45.0 --name=UserSaver
type UserSaver interface {
	SaveUser(
		ctx context.Context,
		nickname string,
		email string,
		passwordHash []byte,
	) (userID string, err error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.45.0 --name=UserProvider
type UserProvider interface {
	User(ctx context.Context, email string) (user models.User, err error)
	IsAdmin(ctx context.Context, userID string) (isAdmin bool, err error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.45.0 --name=AppProvider
type AppProvider interface {
	App(ctx context.Context, appID string) (app models.App, err error)
}

type Authentication struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
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
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     conf.TokenTTL,
	}
}
