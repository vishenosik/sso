package context

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"
)

type appContextKey struct{}

type appContext struct {
	Logger *slog.Logger
}

func NewAppContext(
	log *slog.Logger,
) (*appContext, error) {

	if log == nil {
		return nil, errors.New("shit happens")
	}

	return &appContext{
		Logger: log,
	}, nil
}

func AppContext(ctx context.Context) (*appContext, error) {
	serverCtx, ok := ctx.Value(appContextKey{}).(*appContext)
	if !ok {
		return nil, errors.New("todo: AppContext error")
	}
	return serverCtx, nil
}

func SetAppContext(ctx context.Context, serverCtx *appContext) context.Context {
	return context.WithValue(ctx, appContextKey{}, serverCtx)
}
