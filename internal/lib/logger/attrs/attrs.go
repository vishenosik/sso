package attrs

import "log/slog"

const (
	AttrError     = "err"
	AttrOperation = "operation"
)

func Error(err error) slog.Attr {
	return slog.String(AttrError, err.Error())
}

func Operation(op string) slog.Attr {
	return slog.String(AttrOperation, op)
}
