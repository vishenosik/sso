package context

import (
	"context"
)

type key interface {
	~struct{}
}

type contextValue[keyType key] interface {
	key() keyType
}

func With[keyType key](ctx context.Context, _ctx contextValue[keyType]) context.Context {
	return context.WithValue(ctx, _ctx.key(), _ctx)
}

func From[_type contextValue[keyType], keyType key](ctx context.Context) (*_type, bool) {
	var v _type
	_ctx, ok := ctx.Value(v.key()).(*_type)
	return _ctx, ok
}
