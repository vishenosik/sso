package errors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pkg/errors"
)

const (
	message         = "test message"
	currentMessage  = "currentError"
	expectedMessage = message + ": " + currentMessage
	code            = http.StatusInternalServerError
)

var (
	currentError = errors.New(currentMessage)
	errList      = map[error]int{
		currentError: http.StatusInternalServerError,
	}
)

func newHandlerSuite(t *testing.T) *handler {

	t.Helper()
	t.Parallel()

	return NewHandler(
		WithMessage(message),
		WithCodes(errList),
	)

}

func Test_defaultOpts(t *testing.T) {

	t.Helper()
	t.Parallel()

	hlr := defaultOpts()
	expect := &handler{}

	assert.Equal(t, expect, hlr)

}

func Test_NewHandler(t *testing.T) {
	hlr := newHandlerSuite(t)
	assert.Equal(t, message, hlr.message)

	hlr = NewHandler()
	assert.Empty(t, hlr.codes)
	assert.Empty(t, hlr.message)
}

func Test_Handle(t *testing.T) {

	hlr := newHandlerSuite(t)
	err := hlr.Handle(currentError)

	assert.Equal(t, code, err.Code())
	assert.Equal(t, expectedMessage, err.Error())

	for _, attr := range err.slogAttrs {

		switch attr.Key {

		case "err":
			assert.Equal(t, "err", attr.Key)
			assert.Equal(t, currentMessage, attr.Value.String())

		default:
			t.Error("no needed args provided")
		}

	}

	err = hlr.Handle(errors.New("new err"))
	assert.Equal(t, http.StatusInternalServerError, err.Code())

	err = hlr.Handle(nil)
	assert.Empty(t, err)
}

func Test_WithCodes(t *testing.T) {

	hlr := newHandlerSuite(t)

	var codes ErrorsCodes = map[error]int{
		currentError: http.StatusBadRequest,
	}

	WithCodes(codes)(hlr)
	assert.Equal(t, hlr.codes, codes)
}

func Test_WithMessage(t *testing.T) {

	hlr := newHandlerSuite(t)
	const msg = "msg"

	WithMessage(msg)(hlr)
	assert.Equal(t, hlr.message, msg)
}
