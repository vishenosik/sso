package main

import (
	"fmt"
	"iter"
	_ "net/http/pprof"
	"reflect"
	"slices"

	"github.com/blacksmith-vish/sso/internal/lib/collections"
)

type Employee struct {
	Name string
	Age  int
}

type Employees []Employee

// create a pre-defined list of employees
var employees = Employees{
	{
		Name: "Nikita",
		Age:  23,
	},
	{
		Name: "Vika",
		Age:  52,
	},
	{
		Name: "John",
		Age:  5,
	},
	{
		Name: "Igor",
		Age:  36,
	},
	// {
	// 	Name: "John",
	// 	Age:  5,
	// },
	{
		Name: "Johnson",
		Age:  5,
	},
}

func main() {

	for x, emp := range employees.EmployeeIterator() {
		fmt.Println(x, emp)
	}

	for x := range slices.Values(employees) {
		for k, v := range x.FieldsIter() {
			fmt.Println(k, v)
		}
	}

	names := employees.Names()

	fmt.Println(names)

	_copy := collections.Unique(employees)

	fmt.Println(employees, len(employees), cap(employees))
	fmt.Println(_copy, len(_copy), cap(_copy))
	fmt.Println(collections.HasDuplicates(employees))

}

func Reversed[V any](s []V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !yield(s[i]) {
				return
			}
		}
	}
}

func (empls Employees) EmployeeIterator() iter.Seq2[string, Employee] {
	return func(yield func(string, Employee) bool) {
		for _, empl := range empls {
			if !yield(empl.Name, empl) {
				return
			}
		}
	}
}

func (empls Employees) Names() []string {
	return slices.Collect(empls.NamesIter())
}

func (empls Employees) NamesIter() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, empl := range empls {
			if !yield(empl.Name) {
				return
			}
		}
	}
}

func (empl Employee) FieldsIter() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		_value := reflect.ValueOf(empl)
		_type := reflect.TypeOf(empl)
		for i := range _value.NumField() {
			if !yield(_type.Field(i).Name, _value.Field(i)) {
				return
			}
		}
	}
}
