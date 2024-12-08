package colors

import (
	"github.com/blacksmith-vish/sso/internal/lib/regex"
)

func HighlightNumbers(src string, color uint8) string {
	return regex.NumberRegex.ReplaceAllStringFunc(
		src,
		func(s string) string {
			return colors[color](s)
		},
	)
}
