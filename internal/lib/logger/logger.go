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
	var handler slog.Handler
	switch env {

	case config.EnvProd:
		handler = slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		)

	case config.EnvTest:
		handler = slog.NewJSONHandler(
			io.Discard,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		)

	case config.EnvDev:
		handler = dev.NewHandler(
			os.Stdout,
			slog.LevelDebug,
			dev.WithYamlMarshaller(),
			dev.WithNumbersHighlight(colors.Blue),
			dev.WithKeyWordsHighlight(map[string]colors.ColorCode{
				attrs.AttrError:     colors.Red,
				attrs.AttrOperation: colors.Green,
			}),
		)

	}
	return slog.New(handler)
}
