package dev

import (
	"context"
	"io"
	stdLog "log"
	"log/slog"

	"github.com/blacksmith-vish/sso/internal/lib/colors"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

type Handler struct {
	slog.Handler
	attrs []slog.Attr
	std   *stdLog.Logger
	// syntax highlighter
	high *colors.Higlighter
}

// The signature of the function for setting parameters
type optsFunc func(*Handler)

func NewHandler(
	out io.Writer,
	slogOpts *slog.HandlerOptions,
	opts ...optsFunc,
) *Handler {
	h := &Handler{
		Handler: slog.NewJSONHandler(out, slogOpts),
		std:     stdLog.New(out, "", 0),
		high:    colors.NewHighlighter(),
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func (h *Handler) Handle(_ context.Context, rec slog.Record) error {

	fields := make(map[string]any, rec.NumAttrs())

	rec.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	for _, a := range h.attrs {
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

	attrs = h.high.HighlightNumbers(attrs)

	attrs = h.high.HighlightKeyWords(attrs)

	h.std.Println(
		rec.Time.Format("[15:05:05.000]"),
		level(rec),
		color.CyanString(rec.Message),
		"\n"+attrs,
	)

	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {

	h.attrs = append(h.attrs, attrs...)
	// return h
	return &Handler{
		Handler: h.Handler.WithAttrs(h.attrs),
		std:     h.std,
		attrs:   attrs,
		high:    h.high,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		Handler: h.Handler.WithGroup(name),
		std:     h.std,
		high:    h.high,
	}
}
