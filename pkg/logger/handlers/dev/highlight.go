package dev

import (
	"log/slog"

	"github.com/fatih/color"
	"github.com/vishenosik/sso/pkg/colors"
)

func level(rec slog.Record) string {
	level := rec.Level.String()

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
	return level
}

func WithNumbersHighlight(color colors.ColorCode) optsFunc {
	return func(h *Handler) {
		h.highlight = colors.Modify(h.highlight, colors.WithNumbersHighlight(color))
	}
}

func WithKeyWordsHighlight(keywordsToColors map[string]colors.ColorCode) optsFunc {
	return func(h *Handler) {
		h.highlight = colors.Modify(h.highlight, colors.WithKeyWordsHighlight(keywordsToColors))
	}
}
