package attrs

import (
	// builtin
	"log/slog"
	"time"

	// internal
	time_helper "github.com/blacksmith-vish/sso/internal/lib/helpers/time"
)

const (
	AttrError     = "err"
	AttrOperation = "operation"
	AttrTook      = "took"
)

func Error(err error) slog.Attr {
	return slog.String(AttrError, err.Error())
}

func Operation(op string) slog.Attr {
	return slog.String(AttrOperation, op)
}

func Took(timeStart time.Time) slog.Attr {
	return slog.String(AttrTook, time_helper.FormatWithMeasurementUnit(time.Since(timeStart)))
}
