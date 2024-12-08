package regex

import (
	"regexp"
)

// regex patterns
var (
	// regex numbers matches numbers not in words or followed by anything
	patternNumber = `(^|\W)(\d+(?:\.\d+)?)(|!\w)`
	// patternWord matches words
	patternWord = `\b` + regexp.QuoteMeta("word") + `\b`
	// patternCombined matches numbers not in words or followed by anything
	patternCombined = `(?<!\w)` + patternNumber + `(?!\w)` + `|` + `(?<!\w)` + patternWord + `(?!\w)`
)

var (
	// combinedRegex compiles the combined pattern
	NumberRegex = regexp.MustCompile(patternNumber)
)

// NumberFinder finds all numbers in text that are not included in words or followed by anything
func NumberFinder(text string) []string {
	matches := NumberRegex.FindAllStringSubmatch(text, -1)

	result := make([]string, 0, len(matches))
	for _, match := range matches {
		result = append(result, match[2]) // The number is captured in the second group
	}

	return result
}
