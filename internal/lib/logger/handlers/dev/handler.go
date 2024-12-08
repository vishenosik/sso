package dev

import (
	"context"
	"io"
	stdLog "log"
	"log/slog"
	"regexp"

	"github.com/blacksmith-vish/sso/internal/lib/colors"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

type Handler struct {
	slog.Handler
	std           *stdLog.Logger
	attrs         []slog.Attr
	highlightNums bool
}

// Сигнатура функции для задания параметров
type optsFunc func(*Handler)

func NewHandler(
	out io.Writer,
	slogOpts *slog.HandlerOptions,
	opts ...optsFunc,
) *Handler {
	h := &Handler{
		Handler: slog.NewJSONHandler(out, slogOpts),
		std:     stdLog.New(out, "", 0),
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func (handler *Handler) Handle(_ context.Context, rec slog.Record) error {

	fields := make(map[string]any, rec.NumAttrs())

	rec.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	for _, a := range handler.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var (
		data []byte
		err  error
	)

	if len(fields) > 0 {
		data, err = yaml.Marshal(fields)
		// data, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	attrs := string(data)
	key := "port"

	if handler.highlightNums {
		attrs = colors.HighlightNumbers(attrs, colors.Blue)
	}

	// TODO:
	// 1. Add more fields like service, env, etc.
	// 2. Move to another package for better reusability
	pattern := `\b` + regexp.QuoteMeta(key) + `\b` + `|` + `\b` + regexp.QuoteMeta("op") + `\b`
	attrs = regexp.MustCompile(pattern).ReplaceAllStringFunc(attrs, func(s string) string {
		return color.RedString(s)
	})

	handler.std.Println(
		rec.Time.Format("[15:05:05.000]"),
		level(rec),
		color.CyanString(rec.Message)+"\n",
		color.WhiteString(attrs),
	)

	return nil
}

func (handler *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {

	handler.attrs = append(handler.attrs, attrs...)
	// return handler
	return &Handler{
		Handler:       handler.Handler.WithAttrs(handler.attrs),
		std:           handler.std,
		attrs:         attrs,
		highlightNums: handler.highlightNums,
	}
}

func (handler *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		Handler:       handler.Handler.WithGroup(name),
		std:           handler.std,
		highlightNums: handler.highlightNums,
	}
}
