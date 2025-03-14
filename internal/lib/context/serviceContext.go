package context

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"
)

type serviceContextKey struct{}

type serviceContext struct {
	Logger *slog.Logger
}

func NewServiceContext(
	log *slog.Logger,
) (*serviceContext, error) {

	if log == nil {
		return nil, errors.New("shit happens")
	}

	return &serviceContext{
		Logger: log,
	}, nil
}

func ServiceContext(ctx context.Context) *serviceContext {
	serverCtx, ok := ctx.Value(serviceContextKey{}).(*serviceContext)
	if !ok {
		return nil
	}
	return serverCtx
}

func SetServiceContext(ctx context.Context, serverCtx *serviceContext) context.Context {
	return context.WithValue(ctx, serviceContextKey{}, serverCtx)
}
