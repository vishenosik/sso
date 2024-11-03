package operation

import "github.com/pkg/errors"

func ReturnFailWithError[T any](result T, op string) func(err error) (T, error) {
	return func(err error) (T, error) {
		return result, errors.Wrap(err, op)
	}
}
