package context

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/blacksmith-vish/sso/pkg/colors"
	pkgctx "github.com/blacksmith-vish/sso/pkg/context"
	"github.com/blacksmith-vish/sso/pkg/logger/attrs"
	"github.com/blacksmith-vish/sso/pkg/logger/handlers/dev"
	"github.com/pkg/errors"
)

const (
	EnvDev  = "dev"
	EnvProd = "prod"
	EnvTest = "test"
)

var errAppCtxLoad = errors.New("failed to get app context")

type appContextKey struct{}

type appContext struct {
	Logger *slog.Logger
	Config Config
}

func (ctx *appContext) Key() appContextKey {
	return appContextKey{}
}

func SetupAppCtx() context.Context {
	return WithAppCtx(context.Background())
}

func WithAppCtx(ctx context.Context) context.Context {

	conf := mustLoadEnvConfig()
	log := setupLogger(conf.Env)

	log.Debug("config loaded from env", slog.Any("config", conf))

	return pkgctx.With(ctx, &appContext{
		Logger: log,
		Config: conf,
	})
}

func AppCtx(ctx context.Context) *appContext {
	contx, ok := pkgctx.From[*appContext](ctx)
	if !ok {
		panic(errAppCtxLoad)
	}
	return contx
}

func setupLogger(env string) *slog.Logger {
	var handler slog.Handler
	switch env {

	case EnvProd:
		handler = slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		)

	case EnvTest:
		handler = slog.NewJSONHandler(
			io.Discard,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		)

	case EnvDev:
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

	default:
		handler = slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		)

	}
	return slog.New(handler)
}
