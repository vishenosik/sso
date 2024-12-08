package main

import (
	"fmt"
	"regexp"

	"github.com/fatih/color"
)

// NumberFinder finds all numbers in text that are not included in words
// NumberFinder finds all numbers in text that are not included in words or followed by anything
func NumberFinder(text string) []string {
	// Regex pattern explanation:
	// (^|\W)      - Match the start of the string or any non-word character
	// (           - Start of capturing group
	//   -?        - Optional minus sign for negative numbers
	//   \d+       - One or more digits
	//   (?:\.\d+)? - Optional decimal point followed by one or more digits
	// )           - End of capturing group
	// (?:\W|$)    - Match any non-word character or the end of the string
	pattern := `(^|\W)(-?\d+(?:\.\d+)?)(?:\W|$)`

	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(text, -1)

	result := make([]string, 0, len(matches))
	for _, match := range matches {
		result = append(result, match[2]) // The number is captured in the second group
	}

	return result
}

// Example usage
func ExampleNumberFinder() {
	text := "There are 42 apples and 3.14pies. The code123 is invalid, but 456 is okay."
	numbers := NumberFinder(text)
	fmt.Println("Found numbers:", numbers)
	// Output: Found numbers: [42 3.14 456]
}

func _main() {
	texts := []string{
		"There are 42 apples and 3.14pies.",
		"The code123 is invalid, but 456 is okay.",
		"Temperature is -5.6 degrees Celsius.",
		"The year is 2023, and the price is $19.99.",
		"a16fcc5e-d4de-4cf9-813f-e7ccf36f29d3",
	}

	for _, text := range texts {
		numbers := NumberFinder(text)
		fmt.Printf("Text: %s\nFound numbers: %v\n\n", text, numbers)
	}
}

func main() {
	texts := []string{
		"There are 42 apples and 3.14pies.",
		"The code123 is invalid, but 456 is okay.",
		"Temperature is -5.6 degrees Celsius.",
		"The year is 2023, and the price is $19.99.",
		"a16fcc5e-d4de-4cf9-813f-e7ccf36f29d3",
	}
	for _, text := range texts {
		patternNumber := `(^|\W)(\d+(?:\.\d+)?)(|!\w)`
		text = regexp.MustCompile(patternNumber).ReplaceAllStringFunc(text, func(s string) string {
			return color.BlueString(s)
		})

		fmt.Println(text)
	}
}
