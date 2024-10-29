package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/blacksmith-vish/sso/internal/lib/config"
	"github.com/blacksmith-vish/sso/internal/lib/logger/handlers/dev"
)

func SetupLogger(env string) *slog.Logger {

	devHandler := dev.NewHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)

	switch env {

	case config.EnvDev:
		return slog.New(devHandler)

	case config.EnvProd:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	case config.EnvTest:
		return slog.New(
			slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	}

	return nil
}
