package errors

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
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

type handlerSuite struct {
	hlr     *handler
	recoder *httptest.ResponseRecorder
	log     *slog.Logger
}

func newHandlerSuite(t *testing.T) handlerSuite {

	t.Helper()
	t.Parallel()

	recoder := httptest.NewRecorder()

	return handlerSuite{
		hlr: NewHandler(
			slog.Default(),
			recoder,
			WithMessage(message),
			WithCodes(errList),
		),
		recoder: recoder,
		log: slog.New(
			slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo}),
		),
	}
}

func Test_defaultOpts(t *testing.T) {

	t.Helper()
	t.Parallel()

	recoder := httptest.NewRecorder()
	log := slog.New(
		slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)

	expect := &handler{
		responseWriter: recoder,
		log:            log,
	}

	hlr := defaultOpts(log, recoder)

	assert.Equal(t, expect, hlr)

}

func Test_NewHandler(t *testing.T) {
	suite := newHandlerSuite(t)
	assert.Equal(t, message, suite.hlr.message)

	suite.hlr = NewHandler(
		suite.log,
		suite.recoder,
		nil,
		WithMessage(""),
	)
	assert.Empty(t, suite.hlr.codes)
	assert.Empty(t, suite.hlr.message)
}

func Test_Handle(t *testing.T) {

	suite := newHandlerSuite(t)
	suite.hlr.Handle(currentError)

	assert.Equal(t, code, suite.recoder.Code)

	suite.hlr.Handle(errors.New("new err"))
	assert.Equal(t, code, suite.recoder.Code)

	suite.hlr.Handle(nil)
	assert.Equal(t, code, suite.recoder.Code)
}

func Test_WithCodes(t *testing.T) {

	suite := newHandlerSuite(t)

	var codes ErrorsCodes = map[error]int{
		currentError: http.StatusBadRequest,
	}

	WithCodes(codes)(suite.hlr)
	assert.Equal(t, suite.hlr.codes, codes)
}

func Test_WithMessage(t *testing.T) {

	suite := newHandlerSuite(t)
	const msg = "msg"

	WithMessage(msg)(suite.hlr)
	assert.Equal(t, suite.hlr.message, msg)
}
