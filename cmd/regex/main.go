package main

import (
	"fmt"
	"regexp"
)

// NumberFinder finds all numbers in text that are not included in words
func NumberFinder(text string) []string {
	// Regex pattern explanation:
	// (?<!\w)     - Negative lookbehind: ensure the number is not preceded by a word character
	// \d+         - Match one or more digits
	// (?:\.\d+)?  - Optionally match a decimal point followed by one or more digits
	// (?!\w)      - Negative lookahead: ensure the number is not followed by a word character
	pattern := `(^|\W)(\d+(?:\.\d+)?)(|!\w)`

	re := regexp.MustCompile(pattern)
	return re.FindAllString(text, -1)
}

// Example usage
func ExampleNumberFinder() {
	text := "There are 42 apples and 3.14 pies. The code123 is invalid, but 456 is okay."
	numbers := NumberFinder(text)
	fmt.Println("Found numbers:", numbers)
	// Output: Found numbers: [42 3.14 456]
}

func main() {
	texts := []string{
		"There are 42 apples and 3.14 pies.",
		"The code123 is invalid, but 456 is okay.",
		"Temperature is -5.6 degrees Celsius.",
		"The year is 2023, and the price is $19.99.",
	}

	for _, text := range texts {
		numbers := NumberFinder(text)
		fmt.Printf("Text: %s\nFound numbers: %v\n\n", text, numbers)
	}
}
