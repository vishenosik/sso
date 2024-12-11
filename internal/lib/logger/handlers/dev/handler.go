package dev

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"github.com/blacksmith-vish/sso/internal/lib/colors"
	"github.com/fatih/color"
)

const (
	timeFormat = "15:04:05.000"
)

type Handler struct {
	handler slog.Handler
	writer  io.Writer
	rec     nextFunc
	buf     *bytes.Buffer
	mutex   *sync.Mutex

	outputEmptyAttrs bool
	// syntax highlighter
	highlight *colors.Higlighter

	// marshaller type
	marshalType uint8
}

// The signature of the function for setting parameters
type optsFunc func(*Handler)

func NewHandler(
	writer io.Writer,
	level slog.Level,
	opts ...optsFunc,
) *Handler {

	handlerOptions := &slog.HandlerOptions{
		Level: level,
	}

	buf := &bytes.Buffer{}

	h := &Handler{
		handler: slog.NewJSONHandler(buf, &slog.HandlerOptions{
			Level:       handlerOptions.Level,
			AddSource:   handlerOptions.AddSource,
			ReplaceAttr: suppressDefaultAttrs(handlerOptions.ReplaceAttr),
		}),
		buf:       buf,
		writer:    writer,
		highlight: colors.NewHighlighter(),
		rec:       handlerOptions.ReplaceAttr,
		mutex:     &sync.Mutex{},
	}
	for _, opt := range opts {
		opt(h)
	}

	fmt.Println(h)

	return h
}

func (h *Handler) Handle(ctx context.Context, rec slog.Record) error {

	attrs, err := h.computeAttrs(ctx, rec)
	if err != nil {
		return err
	}

	attrsStr, err := h.marshal(attrs)
	if err != nil {
		return err
	}

	attrsStr = h.highlight.HighlightNumbers(attrsStr)
	attrsStr = h.highlight.HighlightKeyWords(attrsStr)

	output := fmt.Sprintf(
		"[%s] %s: %s\n%s\n",
		rec.Time.Format(timeFormat),
		level(rec),
		color.CyanString(rec.Message),
		attrsStr,
	)

	_, err = io.WriteString(h.writer, output)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{
		handler:     h.handler.WithAttrs(attrs),
		writer:      h.writer,
		highlight:   h.highlight,
		buf:         h.buf,
		rec:         h.rec,
		mutex:       h.mutex,
		marshalType: h.marshalType,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		handler:     h.handler.WithGroup(name),
		writer:      h.writer,
		highlight:   h.highlight,
		buf:         h.buf,
		rec:         h.rec,
		mutex:       h.mutex,
		marshalType: h.marshalType,
	}
}
