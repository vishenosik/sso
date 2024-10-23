package errors

import "log/slog"

type httpError struct {
	err       error
	code      int
	slogAttrs []slog.Attr
}

func (err *httpError) Error() string {
	return err.err.Error()
}

func (err *httpError) Code() int {
	return err.code
}

func (err *httpError) SlogAttrs() []any {
	out := make([]any, len(err.slogAttrs))
	for i := range err.slogAttrs {
		out[i] = err.slogAttrs[i]
	}
	return out
}
