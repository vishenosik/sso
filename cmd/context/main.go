package main

import (
	"context"
	"log/slog"
)

func main() {

	ctx := context.Background()

	// ctx = SetContext(ctx, NewServerContext(logger.SetupLogger("dev")))

	cont := FromContext(ctx)

	cont.Logger = cont.Logger.With(
		slog.String("op", "Authentication_Login_FullMethodName"),
	)

	cont.Logger.Error("shit 1")

	// ctx = SetContext(ctx, NewServerContext(logger.SetupLogger("dev")))

	cont = FromContext(ctx)

	cont.Logger = cont.Logger.With(
		slog.String("op", "sdvsdvsdvsdv"),
		slog.String("op", "aaaaaaaaaaaaa"),
	)

	cont.Logger.Error("shit 2")

}

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
