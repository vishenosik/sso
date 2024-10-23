package errors

import (
	"log/slog"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	errAttr = slog.String("err", currentError.Error())

	attrsList = []any{
		errAttr,
	}
)

func newHttpErrorSuite(t *testing.T) httpError {

	t.Helper()
	t.Parallel()

	return httpError{
		err:  errors.Wrap(currentError, message),
		code: code,
		slogAttrs: []slog.Attr{
			errAttr,
		},
	}

}

func Test_Error(t *testing.T) {
	err := newHttpErrorSuite(t)
	assert.Equal(t, expectedMessage, err.Error())
}

func Test_Code(t *testing.T) {
	err := newHttpErrorSuite(t)
	assert.Equal(t, code, err.Code())
}

func Test_SlogAttrs(t *testing.T) {
	err := newHttpErrorSuite(t)
	assert.Equal(t, attrsList, err.SlogAttrs())
}
