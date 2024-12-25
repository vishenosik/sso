package operation

import (
	"github.com/pkg/errors"
)

func FailResult[Type any](result Type, op string) func(err error) (Type, error) {
	return func(err error) (Type, error) {
		return result, errors.Wrap(err, op)
	}
}
