package context

import (
	"context"
	"log/slog"

	pkgctx "github.com/blacksmith-vish/sso/pkg/context"
)

type appContextKey struct{}

type appContext struct {
	Logger *slog.Logger
}

func (ctx *appContext) Key() appContextKey {
	return appContextKey{}
}

func WithAppCtx(
	ctx context.Context,
	logger *slog.Logger,
) context.Context {
	return pkgctx.With(ctx, &appContext{
		Logger: logger,
	})
}

func AppCtx(ctx context.Context) (*appContext, bool) {
	return pkgctx.From[*appContext](ctx)
}
