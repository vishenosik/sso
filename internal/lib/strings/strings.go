package strings

import (
	"strings"

	"github.com/blacksmith-vish/sso/internal/lib/collections"
)

func ReplaceAllStringFunc(
	src string,
	replacements []string,
	replaceFunc func(string) string,
) string {
	replacements = collections.Unique(replacements)
	for i := range replacements {
		src = strings.ReplaceAll(
			src,
			replacements[i],
			replaceFunc(replacements[i]),
		)
	}
	return src
}
