package attrs

import "log/slog"

func Error(err error) slog.Attr {
	return slog.String("err", err.Error())
}

func Operation(op string) slog.Attr {
	return slog.String("operation", op)
}
