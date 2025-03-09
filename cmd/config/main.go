package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {

	builder := new(strings.Builder)

	Func(builder)

	builder.Write([]byte("func\n"))

	fmt.Println(builder.String())

}

func Func(writer io.Writer) {

	writer.Write([]byte("func"))

}
