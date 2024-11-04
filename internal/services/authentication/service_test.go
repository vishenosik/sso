package authentication

import (
	"io"
	"log/slog"
	"time"

	"github.com/blacksmith-vish/sso/internal/lib/config"
	"github.com/brianvoe/gofakeit/v6"
)

const (
	appSecret      = "secret"
	passDefautlLen = 10
	WrongID        = "invalidx-uuid-xxxx-xxxx-xxxxxxxxxxxx"
)

func suite_newConfig() config.AuthenticationService {
	return config.AuthenticationService{
		TokenTTL: time.Minute,
	}
}

func suite_NewService(
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
) *Authentication {
	return NewService(
		slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})),
		suite_newConfig(),
		userSaver,
		userProvider,
		appProvider,
	)
}

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefautlLen)
}
