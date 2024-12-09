package colors

import (
	"regexp"

	"github.com/blacksmith-vish/sso/internal/lib/regex"
)

type Higlighter struct {
	numbers          *regexp.Regexp
	doNumbers        bool
	numbersColor     ColorCode
	keywords         *regexp.Regexp
	keywordsToColors map[string]ColorCode
}

// The signature of the function for setting parameters
type optsFunc func(*Higlighter)

func NewHighlighter(
	opts ...optsFunc,
) *Higlighter {
	h := &Higlighter{
		numbers: regex.NumberRegex,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func Modify(
	h *Higlighter,
	opts ...optsFunc,
) *Higlighter {
	if h == nil {
		h = NewHighlighter(opts...)
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func (h *Higlighter) HighlightNumbers(src string) string {
	if h.doNumbers {
		return h.numbers.ReplaceAllStringFunc(src, func(s string) string { return colors[h.numbersColor](s) })
	}
	return src
}

func (h *Higlighter) HighlightKeyWords(src string) string {
	if len(h.keywordsToColors) == 0 {
		return src
	}
	return h.keywords.ReplaceAllStringFunc(src, func(s string) string {
		return colors[h.keywordsToColors[s]](s)
	})

}

func WithNumbersHighlight(color ColorCode) optsFunc {
	return func(h *Higlighter) {
		h.doNumbers = true
		h.numbersColor = color
	}
}

func WithKeyWordsHighlight(keywordsToColors map[string]ColorCode) optsFunc {
	return func(h *Higlighter) {
		keywords := make([]string, 0, len(keywordsToColors))
		for key := range keywordsToColors {
			keywords = append(keywords, key)
		}
		h.keywords = regex.KeyWordsCompile(keywords...)
		h.keywordsToColors = keywordsToColors
	}
}
