package context

import (
	"context"
	"os"

	pkgctx "github.com/vishenosik/sso/pkg/context"
)

type signalContextKey struct{}

type signalContext struct {
	Signal os.Signal
}

func (ctx *signalContext) Key() signalContextKey {
	return signalContextKey{}
}

func WithSignalCtx(
	ctx context.Context,
	signal os.Signal,
) context.Context {
	return pkgctx.With(ctx, &signalContext{
		Signal: signal,
	})
}

func SignalCtx(ctx context.Context) (*signalContext, bool) {
	return pkgctx.From[*signalContext](ctx)
}
