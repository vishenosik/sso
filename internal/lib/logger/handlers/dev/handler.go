package dev

import (
	"context"
	"io"
	stdLog "log"
	"log/slog"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

type DevHandler struct {
	slog.Handler
	std   *stdLog.Logger
	attrs []slog.Attr
}

func NewHandler(
	out io.Writer,
	opts *slog.HandlerOptions,
) *DevHandler {
	h := &DevHandler{
		Handler: slog.NewJSONHandler(out, opts),
		std:     stdLog.New(out, "", 0),
	}
	return h
}

func (handler *DevHandler) Handle(_ context.Context, rec slog.Record) error {

	level := rec.Level.String() + ":"

	switch rec.Level {

	case slog.LevelDebug:
		level = color.MagentaString(level)

	case slog.LevelInfo:
		level = color.BlueString(level)

	case slog.LevelWarn:
		level = color.YellowString(level)

	case slog.LevelError:
		level = color.RedString(level)

	}

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

	timeStr := rec.Time.Format("[15:05:05.000]")
	msg := color.CyanString(rec.Message)

	handler.std.Println(
		timeStr,
		level,
		msg,
	)

	attrs := string(data)
	key := "port"

	attrs = strings.ReplaceAll(attrs, key, color.RedString(key))

	handler.std.Println(
		color.WhiteString(attrs),
	)

	return nil
}

func (handler *DevHandler) WithAttrs(attrs []slog.Attr) slog.Handler {

	handler.attrs = append(handler.attrs, attrs...)
	// return handler
	return &DevHandler{
		Handler: handler.Handler.WithAttrs(handler.attrs),
		std:     handler.std,
		attrs:   attrs,
	}
}

func (handler *DevHandler) WithGroup(name string) slog.Handler {
	return &DevHandler{
		Handler: handler.Handler.WithGroup(name),
		std:     handler.std,
	}
}
