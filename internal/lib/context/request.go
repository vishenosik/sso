package context

import (
	"context"
)

type requestContextKey struct{}

type requestContext struct {
	RequestID string
}

func (ctx requestContext) key() requestContextKey {
	return requestContextKey{}
}

func NewRequestContext(
	requestID string,
) (*requestContext, error) {
	return &requestContext{
		RequestID: requestID,
	}, nil
}

func RequestContext(ctx context.Context) (*requestContext, bool) {
	return From[requestContext](ctx)
}
