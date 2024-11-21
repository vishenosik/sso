package context

import (
	"context"
	"log/slog"
)

type contextKey struct{}

type serverContext struct {
	Logger *slog.Logger
}

func NewServerContext(
	log *slog.Logger,
) *serverContext {
	return &serverContext{
		Logger: log,
	}
}

func FromContext(ctx context.Context) *serverContext {
	serverCtx, ok := ctx.Value(contextKey{}).(*serverContext)
	if !ok {
		return nil
	}
	return serverCtx
}

func SetContext(ctx context.Context, serverCtx *serverContext) context.Context {
	return context.WithValue(ctx, contextKey{}, serverCtx)
}
