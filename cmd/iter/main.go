package main

import (
	"fmt"
	"iter"
	"slices"
)

func main() {
	for x, emp := range EmployeeIterator(Employees) {
		fmt.Println(x, emp)
	}

	for x := range slices.Values(Employees) {
		fmt.Println(x)
	}

	// slices.Collect()
}

// Let's define an arbitrary struct
type Employee struct {
	Name   string
	Salary int
}

// create a pre-defined list of employees
var Employees = []Employee{
	{Name: "Elliot", Salary: 4},
	{Name: "Donna", Salary: 5},
}

func EmployeeIterator(e []Employee) iter.Seq2[int, Employee] {
	return func(yield func(int, Employee) bool) {
		for i := 0; i <= len(e)-1; i++ {
			if !yield(i+1, e[i]) {
				return
			}
		}
	}
}

// our iterator function takes in a value and returns a func
// that takes in another func with a signature of `func(int) bool`
func Countdown(v int) func(func(int) bool) {
	// next, we return a callback func which is typically
	// called yield, but names like next could also be
	// applicable
	return func(yield func(int) bool) {
		// we then start a for loop that iterates
		for i := v; i >= 0; i-- {
			// once we've finished looping
			if !yield(i) {
				// we then return and finish our iterations
				return
			}
		}
	}
}

// func main() {
// 	s := []int{1, 2, 3, 4, 5}
// 	// uses the Reversed iterator defined previously
// 	next, stop := iter.Pull(Reversed(s))
// 	defer stop()

// 	for {
// 		v, ok := next()
// 		if !ok {
// 			break
// 		}
// 		fmt.Print(v, " ")
// 	}
// }

// Reversed returns an iterator that loops over a slice in reverse order.
func Reversed[V any](s []V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for i := len(s) - 1; i >= 0; i-- {
			if !yield(s[i]) {
				return
			}
		}
	}
}

func PrintAll[V any](s iter.Seq[V]) {
	for v := range s {
		fmt.Print(v, " ")
	}
	fmt.Println()
}
