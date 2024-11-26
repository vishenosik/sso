package main

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

type Student struct {
	Name string
	Age  int
}

func main() {
	s1 := Student{
		Name: "Sagar",
		Age:  23,
	}

	yamlData, err := yaml.Marshal(s1)

	if err != nil {
		fmt.Printf("Error while Marshaling struct. %v", err)
	}

	fmt.Println(" --- YAML STRUCT---")
	fmt.Println(string(yamlData)) // yamlData will be in bytes. So converting it to string.

	m := make(map[string]any)
	m["name"] = "test"
	m["age"] = 30
	m["err"] = "test-err"

	m["group"] = map[string]any{
		"test": "test-in",
		"err":  "test-err-in",
	}

	m["students"] = []Student{
		{
			Name: "Sagar",
			Age:  23,
		},
		{
			Name: "Sagar",
			Age:  25,
		},
	}

	yamlData, err = yaml.Marshal(m)
	if err != nil {
		fmt.Printf("Error while Marshaling map. %v", err)
	}

	data := string(yamlData)

	re := regexp.MustCompile("[0-9]+")
	nums := re.FindAllString(data, -1)

	fmt.Println(nums)

	slices.SortFunc(nums, func(i, j string) int {
		return len(j) - len(i) // Sort by length in descending order.  If lengths are same, sort lexicographically.  This will ensure the order of numbers in the map is preserved.  For example, "123" will come before "12".  If you want to sort lexicographically, remove this line.  For example, "123" will come before "12".  If you want to sort by length in descending order, remove this line.  For example, "123" will come before "12".  If you want to sort by length in descending order, remove this line.  For example, "123" will come before "12".  If you want to sort by length in descending order, remove this line.  For example, "123" will come before "12".  If you want to sort by
	})

	fmt.Println(nums)

	for _, num := range nums {
		data = strings.ReplaceAll(data, num, color.BlueString(num))
	}
	fmt.Println(" --- YAML MAP---")
	fmt.Println(data) // yamlData will be in bytes. So converting it to string.

}
