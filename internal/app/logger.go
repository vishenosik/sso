package app

import (
	"io"
	"log/slog"
	"os"

	"github.com/blacksmith-vish/sso/pkg/colors"
	"github.com/blacksmith-vish/sso/pkg/logger/attrs"
	"github.com/blacksmith-vish/sso/pkg/logger/handlers/dev"
)

func setupLogger(env string) *slog.Logger {
	var handler slog.Handler
	switch env {

	case envProd:
		handler = slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		)

	case envTest:
		handler = slog.NewJSONHandler(
			io.Discard,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		)

	case envDev:
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
