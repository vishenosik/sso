package dev

import (
	"log/slog"

	"github.com/blacksmith-vish/sso/internal/lib/colors"
	"github.com/fatih/color"
)

func level(rec slog.Record) string {
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
	return level
}

func WithNumbersHighlight(color colors.ColorCode) optsFunc {
	return func(h *Handler) {
		h.high = colors.Modify(h.high, colors.WithNumbersHighlight(color))
	}
}

func WithKeyWordsHighlight(keywordsToColors map[string]colors.ColorCode) optsFunc {
	return func(h *Handler) {
		h.high = colors.Modify(h.high, colors.WithKeyWordsHighlight(keywordsToColors))
	}
}
