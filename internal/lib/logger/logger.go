package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/blacksmith-vish/sso/internal/lib/colors"
	"github.com/blacksmith-vish/sso/internal/lib/config"
	"github.com/blacksmith-vish/sso/internal/lib/logger/attrs"
	"github.com/blacksmith-vish/sso/internal/lib/logger/handlers/dev"
)

func SetupLogger(env string) *slog.Logger {

	handlers := map[string]slog.Handler{
		config.EnvDev: dev.NewHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
			dev.WithNumbersHighlight(colors.Blue),
			dev.WithKeyWordsHighlight(map[string]colors.ColorCode{
				attrs.AttrError:     colors.Red,
				attrs.AttrOperation: colors.Green,
			}),
		),

		config.EnvProd: slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		),

		config.EnvTest: slog.NewJSONHandler(
			io.Discard,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		),
	}

	return slog.New(handlers[env])
}
