package operation

import (
	"log/slog"

	"github.com/blacksmith-vish/sso/internal/lib/logger/attrs"
	"github.com/pkg/errors"
)

func ReturnFailWithError[T any](result T, op string) func(err error) (T, error) {
	return func(err error) (T, error) {
		return result, errors.Wrap(err, op)
	}
}

func FailResultWithAttr[T any](result T, op string) (func(err error) (T, error), slog.Attr) {
	return func(err error) (T, error) {
		return result, errors.Wrap(err, op)
	}, attrs.Operation(op)
}
