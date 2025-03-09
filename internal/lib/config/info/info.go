package info

import (
	"flag"
	"io"
	"os"
)

func init() {
	flag.BoolFunc("config.info", "Show config schema information", configInfo(os.Stdout))
}

type Informator interface {
	Schema() []byte
}

type configComponents = []Informator

var components configComponents

func AddToSchema(in Informator) {
	if len(components) == 0 {
		components = make(configComponents, 0)
	}
	components = append(components, in)
}

func configInfo(
	writer io.Writer,
) func(string) error {
	return func(string) error {
		writer.Write([]byte("\nconfig info:\n\n"))
		for _, component := range components {
			writer.Write(component.Schema())
		}
		os.Exit(0)
		return nil
	}
}
