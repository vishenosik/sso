package errors

import (
	"log/slog"
	"net/http"

	"github.com/pkg/errors"
)

type ErrorsCodes map[error]int

type handler struct {
	message string
	codes   ErrorsCodes
}

// Сигнатура функции для задания параметров
type optsFunc func(*handler)

// Задает опции по умолчанию
func defaultOpts() *handler {
	return &handler{}
}

func NewHandler(opts ...optsFunc) *handler {

	hlr := defaultOpts()

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(hlr)
	}

	return hlr
}

func (hlr *handler) Handle(err error) httpError {

	if err == nil {
		return httpError{}
	}

	Error := httpError{
		err: errors.Wrap(err, hlr.message),
		slogAttrs: []slog.Attr{
			slog.String("err", err.Error()),
		},
	}

	for Err, Code := range hlr.codes {

		if errors.Is(err, Err) {
			Error.code = Code
			return Error
		}
	}

	Error.code = http.StatusInternalServerError
	return Error
}

func WithCodes(codes ErrorsCodes) optsFunc {
	return func(hlr *handler) {
		if hlr != nil {
			hlr.codes = codes
		}
	}
}

func WithMessage(message string) optsFunc {
	return func(hlr *handler) {
		if hlr != nil {
			hlr.message = message
		}
	}
}
