package errors

import (
	"log/slog"
	"net/http"

	"github.com/pkg/errors"
)

type ErrorsCodes map[error]int

type handler struct {
	responseWriter http.ResponseWriter
	log            *slog.Logger
	codes          ErrorsCodes
	message        string
}

// Сигнатура функции для задания параметров
type optsFunc func(*handler)

// Задает опции по умолчанию
func defaultOpts(
	log *slog.Logger,
	w http.ResponseWriter,
) *handler {
	return &handler{
		log:            log,
		responseWriter: w,
	}
}

func NewHandler(
	log *slog.Logger,
	w http.ResponseWriter,
	opts ...optsFunc,
) *handler {

	hlr := defaultOpts(log, w)

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(hlr)
	}

	return hlr
}

func (hlr *handler) Handle(err error) {

	var code int = http.StatusInternalServerError

	if err == nil {
		return
	}

	for Err, Code := range hlr.codes {
		if errors.Is(err, Err) {
			code = Code
			break
		}
	}

	hlr.log.Error(hlr.message, slog.String("err", err.Error()))

	http.Error(
		hlr.responseWriter,
		errors.Wrap(err, hlr.message).Error(),
		code,
	)
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
