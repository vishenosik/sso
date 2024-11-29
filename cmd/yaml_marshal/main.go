package main

import (
	"fmt"
	"regexp"

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
			Age:  2,
		},
	}

	yamlData, err = yaml.Marshal(m)
	if err != nil {
		fmt.Printf("Error while Marshaling map. %v", err)
	}

	data := string(yamlData)

	patternNumber := `[0-9]+`
	data = regexp.MustCompile(patternNumber).ReplaceAllStringFunc(data, func(s string) string {
		return color.BlueString(s)
	})

	fmt.Println(" --- YAML MAP---")
	fmt.Println(data) // yamlData will be in bytes. So converting it to string.

}
