package colors

import (
	"github.com/fatih/color"
)

const (
	Red uint8 = iota
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

type colorFunc func(format string, a ...any) string

var colors = map[uint8]colorFunc{
	Red:     color.RedString,
	Green:   color.GreenString,
	Yellow:  color.YellowString,
	Blue:    color.BlueString,
	Magenta: color.MagentaString,
	Cyan:    color.CyanString,
	White:   color.WhiteString,
}
