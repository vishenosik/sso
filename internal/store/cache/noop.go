package cache

import (
	"context"
	"errors"
	"time"
)

var ErrNoop = errors.New("cache is not available")

type NoopCache struct{}

func NewNoopCache() *NoopCache { return new(NoopCache) }

func (noop *NoopCache) Set(
	ctx context.Context,
	key string,
	value any,
	expiration time.Duration,
) error {
	return ErrNoop
}

func (noop *NoopCache) Get(
	ctx context.Context,
	key string,
) (string, error) {
	return "", ErrNoop
}

func (noop *NoopCache) Delete(
	ctx context.Context,
	key string,
) error {
	return ErrNoop
}
